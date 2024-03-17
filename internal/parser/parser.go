package parser

func Parse_Statement(p ParserInfo) (string, ParserInfo, error) {
	vals, next, err := compose(
		take_while1(is_numeric),
		take_while1(is_whitespace),
		parse_statement,
		parse_newline,
	)(p)

	if err != nil {
		return "", p, err
	}

	got := ""
	for _, val := range vals {
		if str, ok := val.(string); ok {
			got += str
		}
	}

	if next.Start >= len(next.Input) {
		next.Is_Done = true
	}

	return got, next, nil
}

func parse_statement(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return alt(
		parse_rem,
		parse_let,
		parse_if_then,
		parse_return,
		parse_end,
		parse_gosub,
		parse_goto,
		parse_print,
	)(p)
}

func parse_newline(p ParserInfo) ([]interface{}, ParserInfo, error) {
	_, next, _ := take_while(is_whitespace)(p)
	newlines, next, err := take_while1(is_newline)(next)
	if err != nil {
		if next.Start >= len(next.Input) {
			return nil, next, nil // reached end of input which is fine
		}
		return nil, next, generate_error(p, "expected newline", "insert a new line here")
	}
	next.Line += len(newlines)
	next.Col = 1
	return []interface{}{newlines}, next, nil
}

func parse_return(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return consume("RETURN")(p)
}

func parse_end(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return consume("END")(p)
}

func parse_rem(p ParserInfo) ([]interface{}, ParserInfo, error) {
	got, next, err := compose(
		consume("REM"),
		take_while1(is_whitespace),
		take_while(not_is_newline),
	)(p)
	return got, next, err
}

func parse_let(p ParserInfo) ([]interface{}, ParserInfo, error) {

	got, next, err := compose(
		consume("LET"),
		take_while1(is_whitespace),
		take_while1(is_alpha),
		take_while(is_whitespace),
		consume("="),
		take_while(is_whitespace),
		parse_expression,
	)(p)
	return got, next, err
}

func parse_if_then(p ParserInfo) ([]interface{}, ParserInfo, error) {
	got, next, err := compose(
		consume("IF"),
		take_while1(is_whitespace),
		parse_expression,
		take_while(is_whitespace),
		alt(
			consume("<="),
			consume(">="),
			consume("<>"),
			consume(">"),
			consume("<"),
			consume("="),
		),
		take_while(is_whitespace),
		parse_expression,
		take_while1(is_whitespace),
		consume("THEN"),
		take_while1(is_whitespace),
		parse_statement,
	)(p)
	return got, next, err
}

func parse_gosub(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return compose(
		consume("GOSUB"),
		take_while1(is_whitespace),
		parse_expression,
	)(p)
}

func parse_goto(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return compose(
		consume("GOTO"),
		take_while1(is_whitespace),
		parse_expression,
	)(p)
}

func parse_expression(p ParserInfo) ([]interface{}, ParserInfo, error) {
	got, next, err := compose(
		alt_optional(
			consume("+"),
			consume("-")),
		parse_term,
		zero_or_many(compose(
			take_while(is_whitespace),
			alt(
				consume("+"),
				consume("-"),
			),
			take_while(is_whitespace),
			parse_term,
		)),
	)(p)
	return got, next, err
}

func parse_term(p ParserInfo) ([]interface{}, ParserInfo, error) {
	got, next, err := compose(
		parse_factor,
		zero_or_many(
			compose(
				take_while(is_whitespace),
				alt(
					consume("*"),
					consume("/"),
				),
				take_while(is_whitespace),
				parse_factor,
			),
		),
	)(p)
	return got, next, err
}

func parse_factor(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return alt(
		take_while1(is_numeric),
		parse_identifier,
		compose(
			consume("("),
			take_while(is_whitespace),
			parse_expression,
			take_while(is_whitespace),
			consume(")"),
		),
	)(p)
}

func parse_identifier(p ParserInfo) ([]interface{}, ParserInfo, error) {
	got, next, err := take_while1(is_alpha)(p)
	if err == nil {
		if str, ok := got[0].(string); ok {
			if len(str) != 1 {
				return nil, p, generate_error(p, "Invalid Variable Name", "Variable names can only be A-Z")
			} else {
				return got, next, err
			}
		}
	}
	return got, next, err
}

func parse_string(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return compose(
		consume("\""),
		take_while(valid_string_char),
		consume("\""),
	)(p)
}

func parse_expr_list(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return compose(
		alt(
			parse_string,
			parse_expression,
		),
		zero_or_many(
			compose(
				take_while(is_whitespace),
				consume(","),
				take_while(is_whitespace),
				alt(
					parse_string,
					parse_expression,
				),
			),
		),
	)(p)
}

func parse_print(p ParserInfo) ([]interface{}, ParserInfo, error) {
	return compose(
		consume("PRINT"),
		take_while1(is_whitespace),
		parse_expr_list,
	)(p)
}
