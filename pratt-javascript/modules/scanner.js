import utils from "./utils.js";

function create_token(type, value, offset) {
    return {
        type,
        value,
        offset,
    };
}

function end(sc) {
    return sc.offset >= sc.length;
}

function skip(sc) {
    while (utils.is_space(sc.characters[sc.offset])) {
        sc.offset += 1;
    }
}

function next(sc) {
    scanner.offset += 1;
    return sc.characters[sc.offset];
}

function peek(sc) {
    if (sc.end()) {
        return EOF;
    }
    return sc.characters[sc.offset + 1];
}

function peek_next(sc) {
    const offset = sc.offset + 2;
    if (offset >= sc.length) {
        return EOF;
    }
    return sc.characters[offset];
}

export function scan(xs) {
    const scanner = {
        source: xs,
        characters: [...xs],
        length: characters.length,
        offset: 0,
        start: 0,
    };
    const m = Object.create(null);
}
