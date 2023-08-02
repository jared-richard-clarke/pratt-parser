import constants from "./constants.js";
import utils from "./utils.js";

function add_token(self, type, value, column, length) {
    self.tokens.push({ type, value, column, length });
}

function is_end(self) {
    return self.current >= self.length;
}

function next(self) {
    const current = self.current;
    self.current += 1;
    return self.characters[current];
}

function peek(self) {
    if (is_end(self)) {
        return constants.EOF;
    }
    return self.characters[self.current];
}

function peek_next(self) {
    const current = self.current + 1;
    if (current >= self.length) {
        return constants.EOF;
    }
    return self.characters[current];
}

function skip_whitespace(self) {
    while (utils.is_space(peek(self))) {
        self.current += 1;
    }
}

function scan_token(self) {
    const char = next(self);
    if (utils.is_space(char)) {
        return;
    } else if (utils.is_operator(char)) {
        add_token(self, char, null, self.start, 1);
        return;
    } else if (utils.is_paren(char)) {
        add_token(self, char, null, self.start, 1);
        // Check for implied multiplication: (7+11)(11+7), or (7+11)7
        if (utils.is_close_paren(char)) {
            skip_whitespace(self);
            const next_char = peek(self);
            if (utils.is_digit(next_char) || utils.is_open_paren(next_char)) {
                add_token(
                    self,
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
        if (utils.is_zero(char) && utils.is_digit(peek(self))) {
            add_token(
                self,
                constants.ERROR,
                constants.LEADING_ZERO,
                self.start,
                1,
            );
            return;
        }
        while (utils.is_digit(peek(self))) {
            next(self);
        }
        if (utils.is_decimal(peek(self)) && utils.is_digit(peek_next(self))) {
            next(self);
            while (utils.is_digit(peek(self))) {
                next(self);
            }
        }
        const parsed_number = Number.parseFloat(
            self.characters.slice(self.start, self.current).join(""),
        );
        // Check for invalid number.
        if (Number.isNaN(parsed_number)) {
            add_token(
                self,
                constants.ERROR,
                constants.NOT_NUMBER,
                self.start,
                self.current - self.start,
            );
            return;
        }
        add_token(
            self,
            constants.NUMBER,
            parsed_number,
            self.start,
            self.current - self.start,
        );
        // Check for implied multiplication: 7(1 + 2)
        skip_whitespace(self);
        if (utils.is_open_paren(peek(self))) {
            add_token(
                self,
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
                self,
                constants.ERROR,
                constants.MISPLACED_DECIMAL,
                self.start,
                1,
            );
            return;
        }
        add_token(self, constants.ERROR, constants.UNKOWN, self.start, 1);
        return;
    }
}

function run(self) {
    while (!is_end(self)) {
        self.start = self.current;
        scan_token(self);
    }
}

export function scan(text) {
    const spread = [...text];
    const lexer = {
        characters: spread,
        tokens: [],
        length: spread.length,
        start: 0,
        current: 0,
    };
    run(lexer);
    return lexer.tokens;
}
