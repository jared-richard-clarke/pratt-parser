import utils from "./utils.js";

let scanner;

function create_token(type, value, offset) {
    return {
        type,
        value,
        offset,
    };
}

function end() {
    return scanner.offset >= scanner.length;
}

function skip() {
    while (utils.is_space(scanner.characters[scanner.offset])) {
        scanner.offset += 1;
    }
}

function next() {
    scanner.offset += 1;
    return scanner.characters[scanner.offset];
}

function peek() {
    if (scanner.end()) {
        return EOF;
    }
    return scanner.characters[scanner.offset + 1];
}

function peek_next() {
    const offset = scanner.offset + 2;
    if (offset >= scanner.length) {
        return EOF;
    }
    return scanner.characters[offset];
}

const mod = Object.create(null);

mod.init = function (xs) {
    scanner = {
        source: xs,
        characters: [...xs],
        length: characters.length,
        offset: 0,
        start: 0,
    };
    return mod;
};
