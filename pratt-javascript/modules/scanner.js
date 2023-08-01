import { EOF, ERROR, NUMBER } from "./constants.js";
import utils from "./utils.js";

function add_token(self, type, value, offset) {
    self.tokens.push({ type, value, offset });
}

function is_end(self) {
    return self.offset >= self.length;
}

function next(self) {
    self.offset += 1;
    return self.characters[self.offset];
}

function peek(self) {
    if (self.end()) {
        return EOF;
    }
    return self.characters[self.offset + 1];
}

function peek_next(self) {
    const offset = self.offset + 2;
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
        add_token(self, char, null, self.offset);
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
        add_token(self, NUMBER, value, self.start);
        return;
    } else {
        add_token(self, ERROR, value, self.offset);
        return;
    }
}

function run(self) {
    while (!is_end(self)) {
        self.start = self.offset;
        scan_token(self);
    }
}

export function scan(xs) {
    const scanner = {
        characters: [...xs],
        tokens: [],
        length: characters.length,
        offset: 0,
        start: 0,
    };
    run(scanner);
    return scanner.tokens;
}
