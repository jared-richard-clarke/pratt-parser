import constants from "./constants.js";
import utils from "./utils.js";

const lexer = (function () {
    // === private methods ===
    function add_token(type, value, column, length) {
        state.tokens.push({ type, value, column, length });
    }
    function is_end() {
        return state.current > state.end;
    }
    function next() {
        const current = state.current;
        state.current += 1;
        return state.characters[current];
    }
    function peek() {
        if (is_end()) {
            return constants.EOF;
        }
        return state.characters[state.current];
    }
    function peek_next() {
        const current = state.current + 1;
        if (current > state.end) {
            return constants.EOF;
        }
        return state.characters[current];
    }
    function skip_whitespace() {
        while (utils.is_space(peek())) {
            state.current += 1;
        }
    }
    function scan_token() {
        const char = next();
        if (utils.is_space(char)) {
            return;
        } else if (utils.is_operator(char)) {
            add_token(char, null, state.start, 1);
            return;
        } else if (utils.is_paren(char)) {
            add_token(char, null, state.start, 1);
            // Check for implied multiplication: (7+11)(11+7), or (7+11)7
            if (utils.is_close_paren(char)) {
                skip_whitespace();
                const next_char = peek();
                if (
                    utils.is_digit(next_char) || utils.is_open_paren(next_char)
                ) {
                    add_token(
                        constants.IMPLIED_MULTIPLY,
                        null,
                        null,
                        0,
                    );
                }
            }
            return;
        } else if (utils.is_digit(char)) {
            // Check for leading zero error: 07 + 11
            if (utils.is_zero(char) && utils.is_digit(peek())) {
                add_token(
                    constants.ERROR,
                    constants.LEADING_ZERO,
                    state.start,
                    1,
                );
                return;
            }
            while (utils.is_digit(peek())) {
                next();
            }
            if (utils.is_decimal(peek()) && utils.is_digit(peek_next())) {
                next();
                while (utils.is_digit(peek())) {
                    next();
                }
            }
            const parsed_number = Number.parseFloat(
                state.characters.slice(state.start, state.current).join(""),
            );
            // Check for not a number.
            if (Number.isNaN(parsed_number)) {
                add_token(
                    constants.ERROR,
                    constants.NOT_NUMBER,
                    state.start,
                    state.current - state.start,
                );
                return;
            }
            add_token(
                constants.NUMBER,
                parsed_number,
                state.start,
                state.current - state.start,
            );
            // Check for implied multiplication: 7(1 + 2)
            skip_whitespace();
            if (utils.is_open_paren(peek())) {
                add_token(
                    constants.IMPLIED_MULTIPLY,
                    null,
                    null,
                    0,
                );
            }
            return;
        } else {
            // Check for misplaced decimal point.
            if (utils.is_decimal(char)) {
                add_token(
                    constants.ERROR,
                    constants.MISPLACED_DECIMAL,
                    state.start,
                    1,
                );
                return;
            }
            add_token(constants.ERROR, char, state.start, 1);
            return;
        }
    }
    // === internal state ===
    const state = {
        characters: [],
        tokens: [],
        end: 0,
        start: 0,
        current: 0,
    };
    // === public methods ===
    const m = Object.create(null);

    m.set = function (text) {
        const spread = [...text];
        state.characters = spread;
        state.tokens = [];
        state.end = spread.length - 1;
        state.start = 0;
        state.current = 0;
        return m;
    };
    m.run = function () {
        while (!is_end()) {
            state.start = state.current;
            scan_token();
        }
        add_token(constants.EOF, null, state.end + 1, 0);
        return state.tokens;
    };
    return Object.freeze(m);
})();

export default function scan(text) {
    return lexer.set(text).run();
}
