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
    function look_ahead(x) {
        x -= 1;
        return function () {
            const index = state.current + x;
            if (index > state.end) {
                return constants.EOF;
            }
            return state.characters[index];
        };
    }
    const peek = look_ahead(1);
    const peek_next = look_ahead(2);
    const peek_after_next = look_ahead(3);

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
            // Check for empty parentheses: ().
            skip_whitespace();
            if (utils.is_close_paren(peek())) {
                add_token(
                    constants.ERROR,
                    constants.EMPTY_PARENS,
                    state.start,
                    state.current - state.start,
                );
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
            // Exponential notation: 7e11.
            if (utils.is_exponent(peek()) && utils.is_digit(peek_next())) {
                next();
                while (utils.is_digit(peek())) {
                    next();
                }
            }
            // Exponential notation: 7e[+-]11.
            if (
                utils.is_exponent(peek()) &&
                utils.is_plus_minus(peek_next()) &&
                utils.is_digit(peek_after_next())
            ) {
                next();
                next();
                while (utils.is_digit(peek())) {
                    next();
                }
            }
            const number_text = state.characters.slice(
                state.start,
                state.current,
            )
                .join("");
            add_token(
                constants.NUMBER,
                number_text,
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
        } else if (utils.is_decimal(char)) {
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
        } else if (
            (char === "N") && (peek() === "a") && (peek_next() === "N")
        ) {
            next();
            next();
            add_token(
                constants.ERROR,
                constants.NAN,
                state.start,
                state.current - state.start,
            );
            return;
        } else {
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
