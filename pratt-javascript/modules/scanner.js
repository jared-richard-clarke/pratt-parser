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

function scan_token(self) {
    const char = next(self);
    if (utils.is_space(char)) {
        return;
    } else if (utils.is_operator(char) || utils.is_paren(char)) {
        add_token(self, constants.SYMBOL, char, self.start, 1);
        return;
    } else if (utils.is_digit(char)) {
        if (utils.is_zero(char) && utils.is_digit(peek(self))) {
            add_token(self, constants.ERROR, char, self.start, 1);
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
        const value = parseFloat(
            self.characters.slice(self.start, self.current).join(""),
        );
        add_token(
            self,
            constants.NUMBER,
            value,
            self.start,
            self.current - self.start,
        );
        return;
    } else {
        add_token(self, constants.ERROR, char, self.start, 1);
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
    const scanner = {
        characters: spread,
        tokens: [],
        length: spread.length,
        start: 0,
        current: 0,
    };
    run(scanner);
    return scanner.tokens;
}
