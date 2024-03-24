package resp

import (
	"errors"
	"go-rds/pkg/peeker"
	"go-rds/pkg/stack"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

type respParser struct {
	peeker peeker.Peeker

	state stack.Stack

	values []Value
}

type RespParser interface {
	Write(bs []byte) error
	TakeValue() (Value, bool)
}

// 保证peeker里面至少有一个字符
func (parser *respParser) parseConstructBeginFrame() error {
	bs := parser.peeker.Peek()
	parser.peeker.Discard(1)
	switch bs[0] {
	case '+':
		s := &parseSimpleStringState{
			state: parseSimpleStringStateWantingMore,
		}
		parser.state.Push(s)
	case '-':
		s := &parseSimpleErrorState{
			state: parseSimpleErrorStateWantingMore,
		}
		parser.state.Push(s)
	case ':':
		s := &parseIntegerState{
			state: parseIntegerWantingMore,
		}
		parser.state.Push(s)
	case '$':
		s := &parseBulkStringState{
			state: parseBulkStringWantingLength,
		}
		parser.state.Push(s)
	case '*':
		s := &parseArrayState{
			state: parseArrayWantingLength,
		}
		parser.state.Push(s)
	case '_':
		s := &parseNullState{
			state: parseNullWantingMore,
		}
		parser.state.Push(s)
	case '#':
		s := &parseBoolState{
			state: parseBoolWantingMore,
		}
		parser.state.Push(s)
	case '!':
		s := &parseBulkErrorState{
			state: parseBulkErrorWantingLength,
		}
		parser.state.Push(s)
	case '=':
		s := &parseVerbatimStringState{
			state: parseVerbatimStringStateWantingLength,
		}
		parser.state.Push(s)
	case ',':
		s := &parseDoubleState{
			state: stateDoubleWantingMore,
		}
		parser.state.Push(s)
	case '(':
		s := &parseBigNumState{
			state: stateBigNumWantingMore,
		}
		parser.state.Push(s)
	case '%':
		s := &parseMapState{
			state: parseMapWantingLength,
		}
		parser.state.Push(s)
	case '~':
		s := &parseSetState{
			state: parseSetWantingLength,
		}
		parser.state.Push(s)
	case '>':
		s := &parsePushState{
			state: parsePushWantingLength,
		}
		parser.state.Push(s)
	default:
		return errors.New("unsupported message type: " + string(bs[0]))
	}

	return nil
}

func (parser *respParser) parse() error {
	for {
		// 如果state的size为0, 通过判断第一个字符的类型来判断下一个该取的值是什么类型
		if parser.state.Size() == 0 {
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			err := parser.parseConstructBeginFrame()
			if err != nil {
				return err
			}
		}

		stateIface := parser.state.Peek()
		switch state := stateIface.(type) {
		case *parseSimpleStringState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			// 判断是否有\r\n, 可以结束当前的状态
			idx := strings.Index(string(bs), "\r\n")
			if idx == -1 {
				return nil
			}
			data := bs[0:idx]
			parser.peeker.Discard(idx + 2)
			ss, err := ToSimpleString(string(data))
			if err != nil {
				panic(err)
			}
			parser.state.Pop()

			if parser.state.Size() == 0 {
				parser.values = append(parser.values, ss)
			} else {
				parser.state.Peek().(container).PushEntry(ss)
			}
		case *parseSimpleErrorState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			// 判断是否有\r\n, 可以结束当前的状态
			idx := strings.Index(string(bs), "\r\n")
			if idx == -1 {
				return nil
			}
			data := bs[0:idx]
			parser.peeker.Discard(idx + 2)

			// 检查是否是合法的simple error
			strData := string(data)
			spaceIdx := strings.Index(strData, " ")
			if spaceIdx == -1 {
				return errors.New("simple error message contains no space")
			}
			errType := strData[:spaceIdx]
			for _, v := range errType {
				if !unicode.IsUpper(v) {
					return errors.New("simple error type is not uppercase letter")
				}
			}
			errMsg := strData[spaceIdx+1:]
			ss, err := ToSimpleError(errType, errMsg)
			if err != nil {
				panic(err)
			}
			parser.state.Pop()

			if parser.state.Size() == 0 {
				parser.values = append(parser.values, ss)
			} else {
				parser.state.Peek().(container).PushEntry(ss)
			}
		case *parseIntegerState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}

			bs := parser.peeker.Peek()
			// 判断是否有\r\n, 可以结束当前的状态
			idx := strings.Index(string(bs), "\r\n")
			if idx == -1 {
				return nil
			}
			data := bs[0:idx]
			parser.peeker.Discard(idx + 2)

			// 检查data是否合法
			if len(data) == 0 {
				return errors.New("empty integer")
			}
			if len(data) != 1 && data[0] == '0' {
				return errors.New("integer length is 1, but not '0'")
			}
			i64, err := strconv.ParseInt(string(data), 10, 64)
			if err != nil {
				return err
			}
			parser.state.Pop()

			integer := ToInteger(i64)
			if parser.state.Size() == 0 {
				parser.values = append(parser.values, integer)
			} else {
				parser.state.Peek().(container).PushEntry(integer)
			}
		case *parseBulkStringState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			switch state.state {
			case parseBulkStringWantingLength:
				crclIdx := strings.Index(string(bs), "\r\n")
				if crclIdx == -1 {
					return nil
				}
				lengthData := bs[0:crclIdx]
				if len(lengthData) == 0 {
					return errors.New("bulk string length size is 0")
				}
				if len(lengthData) != 1 && lengthData[0] == '0' {
					return errors.New("bulk string length is 1, buf char is not '0'")
				}
				if lengthData[0] == '+' || lengthData[0] == '-' {
					return errors.New("cannot contain + or -")
				}
				i64, err := strconv.ParseInt(string(lengthData), 10, 64)
				if err != nil {
					return err
				}
				parser.peeker.Discard(crclIdx + 2)
				state.state = parseBulkStringWantingData
				state.length = i64
			case parseBulkStringWantingData:
				if parser.peeker.LenRemain() < int(state.length) {
					return nil
				}
				data := bs[0:state.length]
				bulkString := ToBulkString(string(data))
				parser.peeker.Discard(int(state.length) + 2)
				parser.state.Pop()
				if parser.state.Size() == 0 {
					parser.values = append(parser.values, bulkString)
				} else {
					parser.state.Peek().(container).PushEntry(bulkString)
				}
			}
		case *parseArrayState:
			switch state.state {
			case parseArrayWantingLength:
				if parser.peeker.LenRemain() == 0 {
					return nil
				}
				bs := parser.peeker.Peek()
				crclIdx := strings.Index(string(bs), "\r\n")
				if crclIdx == -1 {
					return nil
				}
				parser.peeker.Discard(crclIdx + 2)

				lengthData := bs[:crclIdx]
				if len(lengthData) == 0 {
					return errors.New("array length is required")
				}
				if len(lengthData) != 1 && lengthData[0] == '0' {
					return errors.New("array length is 1, but not '0")
				}
				if lengthData[0] == '+' || lengthData[0] == '-' {
					return errors.New("array length cannot contain + or -")
				}
				length, err := strconv.ParseInt(string(lengthData), 10, 64)
				if err != nil {
					return errors.New("parse int error: " + err.Error())
				}
				state.length = int(length)
				state.state = parseArrayWantingData
			case parseArrayWantingData:
				if len(state.data) == state.length {
					array := ToArray(state.data...)
					parser.state.Pop()
					if parser.state.Size() == 0 {
						parser.values = append(parser.values, array)
					} else {
						parser.state.Peek().(container).PushEntry(array)
					}
				} else {
					if parser.peeker.LenRemain() == 0 {
						return nil
					}

					err := parser.parseConstructBeginFrame()
					if err != nil {
						return err
					}
				}
			}
		case *parseNullState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			crclIdx := strings.Index(string(bs), "\r\n")
			if crclIdx == -1 {
				return nil
			}
			parser.peeker.Discard(crclIdx + 2)

			nullState := ToNull()

			parser.state.Pop()
			if parser.state.Size() == 0 {
				parser.values = append(parser.values, nullState)
			} else {
				parser.state.Peek().(container).PushEntry(nullState)
			}
		case *parseBoolState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			crclIdx := strings.Index(string(bs), "\r\n")
			if crclIdx == -1 {
				return nil
			}
			parser.peeker.Discard(crclIdx + 2)
			bs = bs[:crclIdx]

			if len(bs) != 1 || (bs[0] != 't' && bs[0] != 'f') {
				return errors.New("bool value is invalid: " + string(bs))
			}
			parser.state.Pop()

			boolObj := ToBool(bs[0] == 't')
			if parser.state.Size() == 0 {
				parser.values = append(parser.values, boolObj)
			} else {
				parser.state.Peek().(container).PushEntry(boolObj)
			}
		case *parseBulkErrorState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			switch state.state {
			case parseBulkErrorWantingLength:
				crclIdx := strings.Index(string(bs), "\r\n")
				if crclIdx == -1 {
					return nil
				}
				lengthData := bs[0:crclIdx]
				parser.peeker.Discard(crclIdx + 2)

				if len(lengthData) == 0 {
					return errors.New("bulk error length size is 0")
				}
				if len(lengthData) != 1 && lengthData[0] == '0' {
					return errors.New("bulk error length is not 1, buf char is not '0'")
				}
				if lengthData[0] == '+' || lengthData[0] == '-' {
					return errors.New("cannot contain + or -")
				}
				i64, err := strconv.ParseInt(string(lengthData), 10, 64)
				if err != nil {
					return err
				}
				if i64 == 0 {
					return errors.New("bulk error cannot be empty string")
				}
				state.state = parseBulkErrorWantingData
				state.length = i64
			case parseBulkErrorWantingData:
				if parser.peeker.LenRemain() < int(state.length) {
					return nil
				}
				data := bs[0:state.length]
				parser.peeker.Discard(int(state.length) + 2)

				// 检查是否是合法的simple error
				strData := string(data)
				spaceIdx := strings.Index(strData, " ")
				if spaceIdx == -1 {
					return errors.New("simple error message contains no space")
				}
				errType := strData[:spaceIdx]
				for _, v := range errType {
					if !unicode.IsUpper(v) {
						return errors.New("simple error type is not uppercase letter")
					}
				}
				errMsg := strData[spaceIdx+1:]
				ss, err := ToBulkError(errType, errMsg)
				if err != nil {
					panic(err)
				}

				parser.state.Pop()
				if parser.state.Size() == 0 {
					parser.values = append(parser.values, ss)
				} else {
					parser.state.Peek().(container).PushEntry(ss)
				}
			}
		case *parseBigNumState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			crclIdx := strings.Index(string(bs), "\r\n")
			if crclIdx == -1 {
				return nil
			}
			data := string(bs[:crclIdx])
			parser.peeker.Discard(crclIdx + 2)
			parser.state.Pop()

			for _, v := range data {
				if !unicode.IsDigit(v) {
					return errors.New("big num invalid: " + data)
				}
			}
			f, _, err := big.ParseFloat(string(data), 10, 0, big.AwayFromZero)
			if err != nil {
				return err
			}
			i := big.NewInt(0)
			bigInt, _ := f.Int(i)
			bigNum := ToBigNum(*bigInt)

			if parser.state.Size() == 0 {
				parser.values = append(parser.values, bigNum)
			} else {
				parser.state.Peek().(container).PushEntry(bigNum)
			}
		case *parseDoubleState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}
			bs := parser.peeker.Peek()
			crclIdx := strings.Index(string(bs), "\r\n")
			if crclIdx == -1 {
				return nil
			}
			data := string(bs[:crclIdx])
			parser.peeker.Discard(crclIdx + 2)
			parser.state.Pop()

			var double Double
			if data == "inf" {
				double.IsPositiveInf = true
			} else if data == "-inf" {
				double.IsNegativeInf = true
			} else if data == "nan" {
				double.IsNan = true
			} else {
				f, err := strconv.ParseFloat(data, 64)
				if err != nil {
					return err
				}
				double.Data = f
			}

			if parser.state.Size() == 0 {
				parser.values = append(parser.values, &double)
			} else {
				parser.state.Peek().(container).PushEntry(&double)
			}
		case *parseVerbatimStringState:
			if parser.peeker.LenRemain() == 0 {
				return nil
			}

			bs := parser.peeker.Peek()
			switch state.state {
			case parseVerbatimStringStateWantingLength:
				crclIdx := strings.Index(string(bs), "\r\n")
				if crclIdx == -1 {
					return nil
				}
				data := bs[:crclIdx]
				parser.peeker.Discard(crclIdx + 2)
				if len(data) == 0 {
					return errors.New("verbatim string invalid: " + string(bs))
				}
				if len(data) != 1 && bs[0] == '0' {
					return errors.New("verbatim string invalid: " + string(bs))
				}
				if data[0] == '+' || data[0] == '-' {
					return errors.New("cannot contain + or - char")
				}

				i64, err := strconv.ParseInt(string(data), 10, 64)
				if err != nil {
					return err
				}
				if i64 < 4 {
					return errors.New("verbatim string invalid: " + string(bs))
				}

				state.length = i64
				state.state = parseVerbatimStringStateWantingData
			case parseVerbatimStringStateWantingData:
				if parser.peeker.LenRemain() < int(state.length) {
					return nil
				}
				data := bs[0:state.length]
				parser.peeker.Discard(int(state.length) + 2)

				// 检查是否是合法的simple error
				strData := string(data)
				spaceIdx := strings.Index(strData, ":")
				if spaceIdx == -1 {
					return errors.New("verbatim string error message contains no space")
				}

				// txt:
				if strData[3] != ':' {
					return errors.New("verbatim string is not valid: " + strData)
				}

				parser.state.Pop()

				typ := strData[:3]
				msg := strData[4:]
				vs, err := ToVerbatimString(typ, msg)
				if err != nil {
					panic("unreachable")
				}

				if parser.state.Size() == 0 {
					parser.values = append(parser.values, vs)
				} else {
					parser.state.Peek().(container).PushEntry(vs)
				}
			}
		case *parseMapState:
			switch state.state {
			case parseMapWantingLength:
				if parser.peeker.LenRemain() == 0 {
					return nil
				}
				bs := parser.peeker.Peek()
				crclIdx := strings.Index(string(bs), "\r\n")
				if crclIdx == -1 {
					return nil
				}
				parser.peeker.Discard(crclIdx + 2)

				lengthData := bs[:crclIdx]
				if len(lengthData) == 0 {
					return errors.New("array length is required")
				}
				if len(lengthData) != 1 && lengthData[0] == '0' {
					return errors.New("array length is not 1, but not '0")
				}
				if lengthData[0] == '+' || lengthData[0] == '-' {
					return errors.New("array length cannot contain + or -")
				}
				length, err := strconv.ParseInt(string(lengthData), 10, 64)
				if err != nil {
					return errors.New("parse int error: " + err.Error())
				}
				state.length = int(length)
				state.state = parseMapWantingData
			case parseMapWantingData:
				if len(state.data) == state.length*2 {
					m := map[Value]Value{}
					for i := 0; i < state.length*2; i += 2 {
						m[state.data[i]] = state.data[i+1]
					}
					mapObj := ToMap(m)
					parser.state.Pop()
					if parser.state.Size() == 0 {
						parser.values = append(parser.values, mapObj)
					} else {
						parser.state.Peek().(container).PushEntry(mapObj)
					}
				} else {
					if parser.peeker.LenRemain() == 0 {
						return nil
					}

					err := parser.parseConstructBeginFrame()
					if err != nil {
						return err
					}
				}
			}
		case *parseSetState:
			switch state.state {
			case parseSetWantingLength:
				if parser.peeker.LenRemain() == 0 {
					return nil
				}
				bs := parser.peeker.Peek()
				crclIdx := strings.Index(string(bs), "\r\n")
				if crclIdx == -1 {
					return nil
				}
				parser.peeker.Discard(crclIdx + 2)

				lengthData := bs[:crclIdx]
				if len(lengthData) == 0 {
					return errors.New("set length is required")
				}
				if len(lengthData) != 1 && lengthData[0] == '0' {
					return errors.New("set length is not 1, but not '0")
				}
				if lengthData[0] == '+' || lengthData[0] == '-' {
					return errors.New("set length cannot contain + or -")
				}
				length, err := strconv.ParseInt(string(lengthData), 10, 64)
				if err != nil {
					return errors.New("parse int error: " + err.Error())
				}
				state.length = int(length)
				state.state = parseSetWantingData
			case parseSetWantingData:
				if len(state.data) == state.length {
					set := ToSet(state.data...)
					parser.state.Pop()
					if parser.state.Size() == 0 {
						parser.values = append(parser.values, set)
					} else {
						parser.state.Peek().(container).PushEntry(set)
					}
				} else {
					if parser.peeker.LenRemain() == 0 {
						return nil
					}

					err := parser.parseConstructBeginFrame()
					if err != nil {
						return err
					}
				}
			}
		case *parsePushState:
			switch state.state {
			case parsePushWantingLength:
				if parser.peeker.LenRemain() == 0 {
					return nil
				}
				bs := parser.peeker.Peek()
				crclIdx := strings.Index(string(bs), "\r\n")
				if crclIdx == -1 {
					return nil
				}
				parser.peeker.Discard(crclIdx + 2)

				lengthData := bs[:crclIdx]
				if len(lengthData) == 0 {
					return errors.New("push length is required")
				}
				if len(lengthData) != 1 && lengthData[0] == '0' {
					return errors.New("push length is not 1, but not '0")
				}
				if lengthData[0] == '+' || lengthData[0] == '-' {
					return errors.New("push length cannot contain + or -")
				}
				length, err := strconv.ParseInt(string(lengthData), 10, 64)
				if err != nil {
					return errors.New("parse int error: " + err.Error())
				}
				state.length = int(length)
				state.state = parsePushWantingData
			case parsePushWantingData:
				if len(state.data) == state.length {
					push := ToPush(state.data...)
					parser.state.Pop()
					if parser.state.Size() == 0 {
						parser.values = append(parser.values, push)
					} else {
						parser.state.Peek().(container).PushEntry(push)
					}
				} else {
					if parser.peeker.LenRemain() == 0 {
						return nil
					}

					err := parser.parseConstructBeginFrame()
					if err != nil {
						return err
					}
				}
			}
		default:
			panic("unreachable")
		}
	}
}

func (parser *respParser) Write(bs []byte) error {
	parser.peeker.Write(bs)
	return parser.parse()
}

func (parser *respParser) TakeValue() (Value, bool) {
	if len(parser.values) != 0 {
		ret := parser.values[0]
		parser.values = parser.values[1:]
		return ret, true
	}
	return nil, false
}

func NewRespParser() RespParser {
	parser := respParser{}
	return &parser
}
