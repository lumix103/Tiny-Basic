package parser

import (
	"bytes"
	"fmt"
	"strings"
)

type ParserInfo struct {
	Input          []byte
	Line           int
	Col            int
	Start          int
	Is_Interactive bool
	Is_Done        bool
}

type ParserFn func(ParserInfo) ([]interface{}, ParserInfo, error)

func take_while(accept func(byte) bool) ParserFn {
	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
		foward := len(p.Input)
		Col := 0
		if p.Start < len(p.Input) {
			for idx, ch := range p.Input[p.Start:] {
				if !accept(ch) {
					foward = idx
					break
				}
				Col++
			}
			//fmt.Println("DEBUG TAKE WHILE: " + fmt.Sprintf("%d %d %d", p.Start+foward, p.Start, foward))
			var got string
			if p.Start+foward >= len(p.Input) {
				got = string(p.Input[p.Start:])
				p.Start = len(p.Input)
			} else {
				got = string(p.Input[p.Start : p.Start+foward])
				p.Start = p.Start + foward
			}

			p.Col = p.Col + Col
			arr := []interface{}{}
			arr = append(arr, got)
			return arr, p, nil
		} else {
			got := ""
			p.Start = len(p.Input)
			arr := []interface{}{}
			arr = append(arr, got)
			return arr, p, nil
		}

	}
}

func take_while1(accept func(byte) bool) ParserFn {
	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
		//fmt.Println("DEBUG: " + string(p.Input[p.Start:]) + " DEBUG START: " + fmt.Sprintf("%d", p.Start))
		got, result, _ := take_while(accept)(p)
		v := got[0]
		if str, ok := v.(string); ok {
			if len(str) != 0 {
				arr := []interface{}{}
				arr = append(arr, str)
				return arr, result, nil
			} else {
				return nil, p, generate_error(p, "did not consume anything", "Did you forget something here?")
			}
		} else {
			return nil, p, fmt.Errorf("failed to convert interface{} to string")
		}
	}
}

// func take_until(accept func(byte) bool) ParserFn {
// 	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
// 		foward := len(p.Input)
// 		Col := 0
// 		for idx, ch := range p.Input[p.Start:] {
// 			if accept(ch) {
// 				foward = idx
// 				break
// 			}
// 			Col++
// 		}
// 		got := string(p.Input[p.Start : p.Start+foward])

// 		if len(got) == 0 {
// 			return nil, p, generate_error(p, "take until", "something")
// 		} else {
// 			p.Start = p.Start + foward
// 			p.Col = Col
// 			arr := []interface{}{}
// 			arr = append(arr, got)
// 			return arr, p, nil
// 		}
// 	}
// }

// func take_until_or_end(accept func(byte) bool) ParserFn {
// 	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
// 		foward := len(p.Input)
// 		Col := 0
// 		for idx, ch := range p.Input[p.Start:] {
// 			if accept(ch) {
// 				fmt.Println("Found", idx, string(ch), "END FOUND")
// 				foward = idx
// 				break
// 			}
// 			Col++
// 		}
// 		var got string
// 		if p.Start+foward < len(p.Input) {
// 			got = string(p.Input[p.Start : p.Start+foward])
// 		} else {
// 			got = string(p.Input[p.Start:])
// 		}
// 		if len(got) == 0 {
// 			return nil, p, generate_error(p, "take until", "something")
// 		} else {
// 			p.Start = len(p.Input)
// 			p.Col = Col
// 			arr := []interface{}{}
// 			arr = append(arr, got)
// 			return arr, p, nil
// 		}
// 	}
// }

func consume(symbol string) ParserFn {
	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
		sym := []byte(symbol)
		l := len(sym)
		if p.Start+l <= len(p.Input) && bytes.Equal(p.Input[p.Start:p.Start+l], sym) {
			arr := []interface{}{}
			arr = append(arr, symbol)
			p.Start = p.Start + l
			p.Col = p.Col + l
			return arr, p, nil
		} else {
			return nil, p, generate_error(p, "expected symbol", fmt.Sprintf("Did you mean `%s`?", symbol))
		}
	}
}

func compose(fn ...ParserFn) ParserFn {
	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
		got := []interface{}{}
		next := p
		for _, f := range fn {
			g, n, err := f(next)
			if err != nil {
				return nil, p, err
			} else {
				got = append(got, g...)
				next = n
			}
		}
		return got, next, nil
	}
}

func alt(fn ...ParserFn) ParserFn {
	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
		for _, f := range fn {
			got, remains, err := f(p)
			if err == nil {
				return got, remains, nil
			}
		}
		return nil, p, generate_error(p, "failed to parse a matching statement", "not a valid statement")
	}
}

func alt_optional(fn ...ParserFn) ParserFn {
	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
		for _, f := range fn {
			got, remains, err := f(p)
			if err == nil {
				return got, remains, nil
			}
		}
		return nil, p, nil
	}
}

func zero_or_many(fn ParserFn) ParserFn {
	return func(p ParserInfo) ([]interface{}, ParserInfo, error) {
		got, next, err := fn(p)
		if err != nil {
			return nil, p, nil // We consumed zero
		} else {
			count := 1
			for {
				g, n, err := fn(next)
				if err == nil {
					next = n
					got = append(got, g...)
					count++
				} else {
					//fmt.Println(n.Start)
					//fmt.Println("DEBUG ZERO OR MANY: ", err)
					break
				}
			}
			//fmt.Println("COUNT: ", count)
			return got, next, nil //We consumed as many
		}
	}
}

func is_alpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func is_numeric(ch byte) bool {
	return (ch >= '0' && ch <= '9')
}

// func is_alphanumeric(ch byte) bool {
// 	return is_alpha(ch) || is_numeric(ch)
// }

func is_whitespace(ch byte) bool {
	return (ch == ' ' || ch == '\t')
}

func is_newline(ch byte) bool {
	return (ch == '\n' || ch == '\r')
}

func not_is_newline(ch byte) bool {
	return !is_newline(ch)
}

func valid_string_char(ch byte) bool {
	return ch >= 32 && ch <= 126 && ch != 34
}

func take_untill_newline_or_end(p ParserInfo) string {
	if p.Start >= len(p.Input) {
		return ""
	} else {
		foward := len(p.Input)
		for idx, ch := range p.Input[p.Start:] {
			if is_newline(ch) {
				foward = idx
				break
			}
		}
		if (p.Start + foward) >= len(p.Input) {
			return string(p.Input[p.Start:])
		} else {

			// fmt.Println(p.Start, len(p.Input), foward)

			return string(p.Input[p.Start : p.Start+foward])
		}
	}
}

func backtrack_to_last_newline(p ParserInfo) string {
	var start = p.Start
	if p.Start >= len(p.Input) {
		start = len(p.Input) - 1
	}
	for i := start; i >= 0; i-- {
		if is_newline(p.Input[i]) {
			if i == start {
				return ""
			}
			return string(p.Input[i+1 : start])
		}
	}
	return string(p.Input[:start])
}

// error: <prompt>
// [<ln>: <Col>] |    <Input>
//
//	^^^^^^^
//
// <suggestion>
func generate_error(p ParserInfo, prompt string, suggestion string) error {
	str := take_untill_newline_or_end(p)
	back := backtrack_to_last_newline(p)
	line := back + str
	//fmt.Println("DEBUG| ", str, back, "| END OF DEBUG")
	line_info := fmt.Sprintf("[%d:%d]", p.Line, p.Col)
	return fmt.Errorf("error: %s\n%s\n%s |   '%s'\n%s %s", prompt, strings.Repeat("-", len(prompt)+7), line_info, line, strings.Repeat(" ", +len(line_info)+6+len(back))+strings.Repeat("^", len(str)), suggestion)
}
