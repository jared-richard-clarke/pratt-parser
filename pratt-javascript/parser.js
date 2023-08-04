import constants from "./modules/constants.js";
import scan from "./modules/lexer.js";
import utils from "./modules/utils.js";

const parser = (function () {
    // parser internal state
    const state = {
        source: [],
        index: 0,
        end: 0,
    };

    const table = (function () {
        const registry = {
            prefix: {},
            infix: {},
            bind: {},
            prebind: {},
        };
        function register(type, parser) {
            registry.prefix[type] = parser;
            registry.bind[type] = 0;
        }
        function register_unary(bp, types, parser) {
            types.forEach((type) => {
                registry.prefix[type] = parser;
                registry.prebind[type] = bp;
            });
        }
        function register_binary(bp, types, parser) {
            types.forEach((type) => {
                registry.infix[type] = parser;
                registry.bind[type] = bp;
            });
        }
        register(constants.EOF, parse_error);
        register(constants.ERROR, parse_error);
        register(constants.NUMBER, parse_literal);
        register(constants.OPEN_PAREN, parse_grouping);
        register_binary(
            10,
            [constants.ADD, constants.SUBTRACT],
            parse_binary(true),
        );
        register_binary(
            20,
            [
                constants.MULTIPLY,
                constants.MULTIPLY_ALT,
                constants.DIVIDE,
                constants.DIVIDE_ALT,
            ],
            parse_binary(true),
        );
        register_binary(30, [constants.EXPONENT], parse_binary(false));
        register_binary(40, [constants.IMPLIED_MULTIPLY], parse_binary(true));
        register_unary(50, [constants.ADD, constants.SUBTRACT], parse_unary);

        const m = Object.create(null);
        m.get = function (category, type) {
            return registry[category][type];
        };
        return Object.freeze(m);
    })();

    function parse_expression(rbp) {
        const token = next();
        const prefix = table.get("prefix", token.type);
        let [left, error] = prefix(token);
        if (error !== null) {
            return [null, error];
        }
        while (rbp < table.get("bind", peek())) {
            const token = next();
            const infix = table.get("infix", token.type);
            [left, error] = infix(left, token);
            if (error !== null) {
                return [null, error];
            }
        }
        return [left, null];
    }

    function parse_error(token) {
        return [null, token];
    }
    function parse_literal(token) {
        return [token.value, null];
    }

    function parse_unary(token) {
        const bind = table.get("prebind", token.type);
        const [x, error] = parse_expression(bind);
        if (error !== null) {
            return [null, error];
        }
        const operation = utils.unary_operation[token.type];
        return [operation(x), null];
    }

    function parse_binary(left) {
        return function (x, token) {
            const bind = table.get("bind", token.type);
            const [y, error] = parse_expression(left ? bind : bind - 1);
            if (error !== null) {
                return [null, error];
            }
            const operation = utils.binary_operation[token.type];
            return [operation(x, y), null];
        };
    }

    function parse_grouping(token) {
        const start_token = token;
        const [x, error] = parse_expression(0);
        if (error !== null) {
            return [null, error];
        }
        if (!match(constants.CLOSE_PAREN)) {
            error.message = start_token;
            return [null, error];
        }
        next();
        return [x, null];
    }
    function next() {
        if (state.index >= state.end) {
            return state.source[state.end];
        }
        const token = state.source[state.index];
        state.index += 1;
        return token;
    }
    function peek() {
        return state.source[state.index].type;
    }
    function match(expect) {
        return peek() === expect;
    }
    // === public methods ===
    const m = Object.create(null);
    m.set = function (text) {
        const tokens = scan(text);
        state.source = tokens;
        state.index = 0;
        state.end = tokens.length - 1;
        return m;
    };
    m.run = function () {
        const [x, error] = parse_expression(0);
        if (error !== null) {
            return [null, error];
        }
        if (state.index < state.end) {
            const error = state.tokens[state.index];
            error.message = constants.INCOMPLETE_EXPRESSION;
            return [null, error];
        }
        return [x, null];
    };
    return Object.freeze(m);
})();

export function parse(text) {
    return parser.set(text).run();
}
