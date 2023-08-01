import { EOF } from "./constants.js";
import utils from "./utils.js";

function add_token(self, value, position, length) {
    self.tokens.push({ value, position, length });
}

function is_end(self) {
    return self.offset >= self.length;
}

function next(self) {
    const offset = self.offset;
    self.offset += 1;
    return self.characters[offset];
}

function peek(self) {
    if (is_end(self)) {
        return EOF;
    }
    return self.characters[self.offset];
}

function peek_next(self) {
    const offset = self.offset + 1;
    if (offset >= self.length) {
        return EOF;
    }
    return self.characters[offset];
}

function scan_token(self) {
    const char = next(self);
    if (utils.is_space(char)) {
        return;
    } else if (utils.is_operator(char) || utils.is_paren(char)) {
        add_token(self, char, self.start, 1);
        return;
    } else if (utils.is_digit(char)) {
        while (utils.is_digit(peek(self))) {
            next(self);
        }
        if (utils.is_decimal(peek(self)) && utils.is_digit(peek_next(self))) {
            next(self);
            while (utils.is_digit(peek(self))) {
                next(self);
            }
        }
        const text = self.characters.slice(self.start, self.offset).join("");
        const value = Number(text);
        add_token(self, value, self.start, self.offset - self.start);
        return;
    } else {
        add_token(self, char, self.start, 1);
        return;
    }
}

function run(self) {
    while (!is_end(self)) {
        self.start = self.offset;
        scan_token(self);
    }
}

export function scan(text) {
    const spread = [...text];
    const scanner = {
        characters: spread,
        tokens: [],
        length: spread.length,
        offset: 0,
        start: 0,
    };
    run(scanner);
    return scanner.tokens;
}
