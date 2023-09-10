import constants from "./modules/constants.js";
import { scan } from "./modules/lexer.js";
import utils from "./modules/utils.js";
import encoders from "./modules/big-math/encoders.js";

// parse(string) -> [string, null] | [null, [token]]
//     where token = { type, value, message, column, length }
//
// Scans, parses, and evaluates a string as an arithmetic expression.
// "parse" does not build an abstract syntax tree. Instead it evaluates the expression as it parses.
// Returns a two part array — "[string, null]" if successful, "[null, [token]]" if unsuccessful.
// "string" is an evaluated arithmetic expression, and "[token]" is an array of error tokens
// both locating and describing errors within the input string.
export const parse = (function () {
    // === parse: state ===
    // Tracks parser's internal state.
    //
    // > state.set(string)
    //   Resets the parser with its new input.
    //
    // > state.next() -> token
    //   Moves the parser to the next token within the source array.
    //   Returns the previous token.
    //
    // > state.peek() -> token.type
    //   Returns the type of the next token without consuming the token.
    //
    // > state.match(string) -> boolean
    //   Checks if the proceeding token type matches the expected token type.
    //
    // > state.consumed() -> boolean
    //   Checks if the parser has consumed all its input.
    //
    // > state.length() -> number
    //   Returns the length of the internal token array.
    //
    // > state.flush(token) -> [token] where token.type = "error"
    //   Collects all the error tokens into a single array and
    //   returns it to the caller.
    const state = (function () {
        const tokens = {
            source: [],
            length: 0,
            index: 0,
            end: 0,
        };
        function set(source) {
            tokens.source = source;
            tokens.length = source.length;
            tokens.index = 0;
            tokens.end = source.length - 1;
        }
        function next() {
            if (tokens.index >= tokens.end) {
                return tokens.source[tokens.end];
            }
            const token = tokens.source[tokens.index];
            tokens.index += 1;
            return token;
        }
        function peek() {
            return tokens.source[tokens.index].type;
        }
        function match(expect) {
            return peek() === expect;
        }
        function consumed() {
            return tokens.index >= tokens.end;
        }
        function length() {
            return tokens.length;
        }
        function flush(error) {
            const errors = [error];
            while (tokens.index < tokens.end) {
                const token = next();
                if (token.type === "error") {
                    errors.push(token);
                }
            }
            return errors;
        }

        return Object.freeze({
            set,
            next,
            peek,
            match,
            consumed,
            length,
            flush,
        });
    })();

    // === parse: parsers ===
    // Top down operator precedence parsing, as imagined by Vaughan Pratt,
    // combines lexical semantics with functions. Each lexeme is assigned a
    // function — its semantic code. To parse a string of lexemes is to execute
    // the semantic code of each lexeme in turn from left to right.
    //
    // There are two types of semantic code:
    // 1. prefix: a lexeme function without a left expression.
    // 2. infix: a lexeme function with a left expression.
    //
    // This semantic code forms the parsers internal to "parse".

    // The engine of Pratt's technique, "parse_expression" drives the parser,
    // calling the semantic code of each lexeme in turn from left to right.
    // For every level of precedence — dictated by binding power — there is a call
    // to "parse_expression" either through the "prefix" parser or "infix" parser
    // of the associated lexeme. The resolution of "parse_expression" is to return
    // either an evaluated expression or an array of error tokens.
    function parse_expression(rbp) {
        const token = state.next();
        const [prefix, ok] = table.get_parser("prefix", token.type);
        if (!ok) {
            token.message += constants.NO_PREFIX;
            return [null, token];
        }
        let [x, error] = prefix(token);
        if (error !== null) {
            return [null, error];
        }
        while (rbp < table.get_binding("bind", state.peek())) {
            const token = state.next();
            const [infix, ok] = table.get_parser("infix", token.type);
            if (!ok) {
                token.message += constants.NO_INFIX;
                return [null, token];
            }
            [x, error] = infix(x, token);
            if (error !== null) {
                return [null, error];
            }
        }
        return [x, null];
    }

    // The "eof" token marks the end of a token array. Calling code
    // on this token means an error. "parse_eof" resolves these errors.
    function parse_eof(token) {
        if (state.length() === 1) {
            // If the expression is empty, then the error spans it.
            token.length = token.column;
            token.column = 0;
            token.message += constants.EMPTY_EXPRESSION;
            return [null, token];
        }
        token.message += constants.INCOMPLETE_EXPRESSION;
        return [null, token];
    }

    // Resolves any error token called in a prefix position.
    function parse_unary_error(token) {
        return [null, token];
    }
    // Resolves any error token called in an infix position.
    function parse_binary_error(x, token) {
        x = null;
        return [x, token];
    }
    // Parses number expressions. Transforms string value into
    // a big float object for evaluation.
    function parse_number(token) {
        const number = encoders.decode(token.value);
        return [number, null];
    }
    // Parses unary expressions. Parses expression, and, if successful,
    // calls the associated unary operation on that expression.
    function parse_unary(token) {
        const bind = table.get_binding("prebind", token.type);
        const [x, error] = parse_expression(bind);
        if (error !== null) {
            return [null, error];
        }
        const operation = utils.unary_operation[token.type];
        return [operation(x), null];
    }

    // Parses binary expressions. Takes a left expression and parses
    // right expression, and, if successful, calls the associated
    // binary operation on both the right and left expressions.
    function parse_binary(left) {
        return function (x, token) {
            const bind = table.get_binding("bind", token.type);
            const [y, error] = parse_expression(left ? bind : bind - 1);
            if (error !== null) {
                return [null, error];
            }
            const operation = utils.binary_operation[token.type];
            const value = operation(x, y);
            if (typeof value === "string") {
                token.message += value;
                return [null, token];
            }
            return [value, null];
        };
    }
    // Parses binary expressions that associate left.
    const parse_left = parse_binary(true);
    // Parses binary expressions that associate right.
    const parse_right = parse_binary(false);

    // Parses expressions grouped within parentheses. If an expression is parsed successfully,
    // following an open parenthesis, "parse_grouping" checks for a matching closed parenthesis.
    function parse_grouping(token) {
        const [x, error] = parse_expression(0);
        if (error !== null) {
            return [null, error];
        }
        if (!state.match(constants.CLOSE_PAREN)) {
            token.message += constants.MISMATCHED_PAREN;
            return [null, token];
        }
        state.next();
        return [x, null];
    }

    // === parse: table ===
    // Maps all the parsers and bindings to their associated lexemes.
    //
    // > table.get_parser(string) -> [parser, true] | [null, false]
    //   If the parser exists for the associated lexeme, returns the
    //   parser alongside boolean true. Otherwise returns null alongside
    //   boolean false.
    //
    // > table.get_binding(string) -> number
    //   Returns the binding power of the associated lexeme.
    const table = (function () {
        // === table: registry ===
        const registry = {
            prefix: {},
            infix: {},
            bind: {},
            prebind: {},
        };
        // Helper functions map the parsers and bindings to their associated lexemes
        // within the lookup table's internal registry.
        function register(bind, type, parser) {
            registry.prefix[type] = parser;
            registry.bind[type] = bind;
        }
        function register_unary(bind, types, parser) {
            types.forEach((type) => {
                registry.prefix[type] = parser;
                registry.prebind[type] = bind;
            });
        }
        function register_binary(bind, types, parser) {
            types.forEach((type) => {
                registry.infix[type] = parser;
                registry.bind[type] = bind;
            });
        }
        // Building the lookup table.
        register(0, constants.EOF, parse_eof);
        register(0, constants.NUMBER, parse_number);
        register(0, constants.OPEN_PAREN, parse_grouping);
        register_binary(
            10,
            [constants.ADD, constants.SUBTRACT, constants.SUBTRACT_ALT],
            parse_left,
        );
        register_binary(
            20,
            [
                constants.MULTIPLY,
                constants.MULTIPLY_ALT,
                constants.DIVIDE,
                constants.DIVIDE_ALT,
            ],
            parse_left,
        );
        register_binary(30, [constants.EXPONENT], parse_right);
        register_binary(40, [constants.IMPLIED_MULTIPLY], parse_left);
        register_unary(
            50,
            [constants.ADD, constants.SUBTRACT, constants.SUBTRACT_ALT],
            parse_unary,
        );
        register(60, constants.ERROR, parse_unary_error);
        register_binary(60, [constants.ERROR], parse_binary_error);

        function get_parser(category, type) {
            const parser = registry[category][type];
            return parser === undefined ? [null, false] : [parser, true];
        }
        function get_binding(category, type) {
            return registry[category][type];
        }

        return Object.freeze({
            get_parser,
            get_binding,
        });
    })();

    // === parser ===
    return function (text) {
        // Transform text into tokens and set internal state.
        const tokens = scan(text);
        state.set(tokens);
        // Parse expression.
        const [x, error] = parse_expression(0);
        if (error !== null) {
            const errors = state.flush(error);
            return [null, errors];
        }
        // Check for unused tokens.
        if (!state.consumed()) {
            const token = state.next();
            if (utils.is_close_paren(token.type)) {
                token.message += constants.MISMATCHED_PAREN;
            } else {
                token.message += constants.INCOMPLETE_EXPRESSION;
            }
            const errors = state.flush(token);
            return [null, errors];
        }
        // Re-encode big-float object as string.
        // Use scientific notation for exceedingly large or small numbers.
        if (utils.is_exceeding(x)) {
            return [encoders.encode_scientific(x), null];
        }
        return [encoders.encode(x), null];
    };
})();

// format(parser, string) -> string
// Applies parser to a string and formats result.
export function format(parser, text) {
    const [success, errors] = parser(text);
    if (success !== null) {
        return success;
    }

    const space = constants.WHITE_SPACE;
    const linefeed = constants.LINEFEED;
    const period = "." + space;
    const caret = "^";
    const expression = text.replace(/\s/g, space) + linefeed;

    let arrows = errors
        .map((error, index, array) => {
            let offset = 0;
            const column = error.column;
            if (index > 0) {
                const prev = array[index - 1];
                offset = prev.column + prev.length;
            }
            return space.repeat(column - offset) + caret;
        });
    arrows.push(linefeed);
    arrows = arrows.join("");

    let messages = errors
        .map((error, index) => {
            const count = String(index + 1) + period;
            return count + error.message;
        });
    messages.push(linefeed);
    messages = messages.join("");
    return expression + arrows + messages;
}
