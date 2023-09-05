import constants from "./constants.js";
import utils from "./utils.js";

// === lexer ===
// A lexer object that transforms a string into an array of token objects.
// Errors have dedicated error tokens, so failure is impossible.
//
// > lexer.input(string) -> lexer:
//   Inputs a string and sets the lexer state. Transforms string into an iterable array of characters.
//   Properties `tokens`, `end`, `start`, and `current` provide the machinery to navigate
//   the characters array.
//   Returns lexer object to the caller to allow method chaining.
//
// > lexer.run() -> [token]:
//   Drives the lexer, transforming a character array into a token array. Characters are one-letter
//   strings whereas tokens are objects containing lexemes, related descriptions, and positional information.
const lexer = (function () {
    // === lexer: state ===
    // Tracks the lexer within the character array.
    //
    // > state.set(string)
    //   Resets the lexer with new input.
    //
    // > state.lexeme() -> string
    //   Returns a lexeme by slicing into the characters array.
    //
    // > state.lexeme_start() -> number
    //   Returns the starting position of the current potential lexeme.
    //
    // > state.lexeme_length() -> number
    //   Returns the length of the current potential lexeme.
    //
    // > state.end() -> number
    //   Returns the final index of the internal characters array.
    //
    // > state.add_token(token)
    //   Appends a new token to the internal tokens array.
    //
    // > state.consumed() -> boolean
    //   Checks if the lexer has consumed all its input.
    //
    // > state.next() -> character
    //   Moves the lexer to the next character within the characters array.
    //   Returns the character to the caller for potential processing.
    //
    // > state.reset()
    //   Moves the lexer to the start of the next potential lexeme.
    //
    // > state.peek() -> string
    //   Look ahead one character.
    //
    // > state.peek_next() -> string
    //   Look ahead two characters.
    //
    // > state.peek_after_next() -> string
    //   Look ahead three characters.
    //
    // > state.skip_whitespace()
    //   Skips whitespace characters while moving the lexer forward.
    //
    // > state.tokens() -> [tokens]
    //   Returns the internal tokens array to the caller.
    const state = (function () {
        const internal = {
            characters: [],
            tokens: [],
            end: 0,
            start: 0,
            current: 0,
        };
        function set(text) {
            const spread = [...text];
            internal.characters = spread;
            internal.tokens = [];
            internal.end = spread.length - 1;
            internal.start = 0;
            internal.current = 0;
        }
        function lexeme() {
            return internal.characters
                .slice(internal.start, internal.current)
                .join("");
        }
        function lexeme_start() {
            return internal.start;
        }
        function lexeme_length() {
            return internal.current - internal.start;
        }
        function end() {
            return internal.end;
        }
        function add_token(type, value, message, column, length) {
            internal.tokens.push({ type, value, message, column, length });
        }
        function consumed() {
            return internal.current > internal.end;
        }
        function next() {
            const current = internal.current;
            internal.current += 1;
            return internal.characters[current];
        }
        function reset() {
            internal.start = internal.current;
        }

        // Factory function produces methods that look ahead by a set number
        // of characters into the character array without the lexer consuming them.
        function look_ahead(x) {
            x -= 1;
            return function () {
                const index = internal.current + x;
                if (index > internal.end) {
                    return constants.EOF;
                }
                return internal.characters[index];
            };
        }
        const peek = look_ahead(1);
        const peek_next = look_ahead(2);
        const peek_after_next = look_ahead(3);

        function skip_whitespace() {
            while (utils.is_space(peek())) {
                internal.current += 1;
            }
        }
        function tokens() {
            return internal.tokens;
        }
        return Object.freeze({
            set,
            lexeme,
            lexeme_start,
            lexeme_length,
            end,
            add_token,
            consumed,
            next,
            reset,
            look_ahead,
            peek,
            peek_next,
            peek_after_next,
            skip_whitespace,
            tokens,
        });
    })();

    // Consumes a lexeme from the characters array, builds a token object,
    // and appends it to the tokens array.
    function scan_token() {
        state.reset();
        const char = state.next();
        if (utils.is_space(char)) {
            return;
        } else if (utils.is_operator(char)) {
            state.add_token(char, null, "", state.lexeme_start(), 1);
            return;
        } else if (utils.is_paren(char)) {
            state.add_token(char, null, "", state.lexeme_start(), 1);
            // Check for implied multiplication: (7+11)(11+7), or (7+11)7
            if (utils.is_close_paren(char)) {
                state.skip_whitespace();
                const next_char = state.peek();
                if (
                    utils.is_digit(next_char) ||
                    utils.is_open_paren(next_char)
                ) {
                    state.add_token(
                        constants.IMPLIED_MULTIPLY,
                        null,
                        "",
                        null,
                        0,
                    );
                }
            }
            // Check for empty parentheses: ().
            state.skip_whitespace();
            if (utils.is_close_paren(state.peek())) {
                state.add_token(
                    constants.ERROR,
                    null,
                    constants.EMPTY_PARENS,
                    state.lexeme_start(),
                    state.lexeme_length(),
                );
            }
            return;
        } else if (utils.is_digit(char)) {
            // Check for leading zero error: 07 + 11
            if (utils.is_zero(char) && utils.is_digit(state.peek())) {
                state.add_token(
                    constants.ERROR,
                    constants.ZERO,
                    constants.LEADING_ZERO,
                    state.lexeme_start(),
                    1,
                );
                return;
            }
            while (utils.is_digit(state.peek())) {
                state.next();
            }
            if (
                utils.is_decimal(state.peek()) &&
                utils.is_digit(state.peek_next())
            ) {
                state.next();
                while (utils.is_digit(state.peek())) {
                    state.next();
                }
            }
            // Exponential notation: 7e11.
            if (
                utils.is_exponent_suffix(state.peek()) &&
                utils.is_digit(state.peek_next())
            ) {
                state.next();
                while (utils.is_digit(state.peek())) {
                    state.next();
                }
            }
            // Exponential notation: 7e[+-]11.
            if (
                utils.is_exponent_suffix(state.peek()) &&
                utils.is_plus_minus(state.peek_next()) &&
                utils.is_digit(state.peek_after_next())
            ) {
                state.next();
                state.next();
                while (utils.is_digit(state.peek())) {
                    state.next();
                }
            }
            state.add_token(
                constants.NUMBER,
                state.lexeme(),
                "",
                state.lexeme_start(),
                state.lexeme_length(),
            );
            // Check for implied multiplication: 7(1 + 2)
            state.skip_whitespace();
            if (utils.is_open_paren(state.peek())) {
                state.add_token(constants.IMPLIED_MULTIPLY, null, "", null, 0);
            }
            return;
        } else if (utils.is_decimal(char)) {
            // Check for misplaced decimal point.
            state.add_token(
                constants.ERROR,
                constants.DECIMAL_POINT,
                constants.MISPLACED_DECIMAL,
                state.lexeme_start(),
                1,
            );
            return;
        } else if (utils.is_ascii_letter(char)) {
            state.next();
            while (utils.is_ascii_letter(state.peek())) {
                state.next();
            }
            const lexeme = state.lexeme();
            // Check for NaN: not a number.
            if (utils.is_nan(lexeme)) {
                state.add_token(
                    constants.ERROR,
                    constants.NAN,
                    constants.NOT_NUMBER,
                    state.lexeme_start(),
                    state.lexeme_length(),
                );
                return;
            }
            // Check for misplaced exponent suffix: "e" or "E".
            if (utils.is_exponent_suffix(lexeme)) {
                state.add_token(
                    constants.ERROR,
                    char,
                    constants.MISPLACED_EXPONENT,
                    state.lexeme_start(),
                    1,
                );
                return;
            }
            state.add_token(
                constants.ERROR,
                lexeme,
                constants.UNKNOWN,
                state.lexeme_start(),
                state.lexeme_length(),
            );
            return;
        } else {
            state.add_token(
                constants.ERROR,
                char,
                constants.UNKNOWN,
                state.lexeme_start(),
                1,
            );
            return;
        }
    }

    // === lexer: public methods ===
    const methods = Object.create(null);

    methods.input = function (text) {
        state.set(text);
        return methods;
    };
    methods.run = function () {
        while (!state.consumed()) {
            scan_token();
        }
        state.add_token(constants.EOF, null, "", state.end() + 1, 0);
        return state.tokens();
    };
    return Object.freeze(methods);
})();

// === scan ===
// scan(string) -> [token]
//     where token = { type, value, message, column, length }
//
// Wraps the lexer object in a simpler one-function API.
// Inputs text to the lexer, runs the lexer, and returns the lexer output.
export default function scan(text) {
    return lexer.input(text).run();
}
