export default Object.freeze({
    // spaces
    WHITE_SPACE: " ",
    TAB: "\t",
    LINEFEED: "\n",
    CARRIAGE_RETURN: "\r",
    VERTICAL_TAB: "\v",
    FORM_FEED: "\f",
    // general symbols
    DECIMAL_POINT: ".",
    OPEN_PAREN: "(",
    CLOSE_PAREN: ")",
    // end-of-input flag
    EOF: "eof",
    // digits
    ZERO: "0",
    ONE: "1",
    TWO: "2",
    THREE: "3",
    FOUR: "4",
    FIVE: "5",
    SIX: "6",
    SEVEN: "7",
    EIGHT: "8",
    NINE: "9",
    // operators
    ADD: "+",
    SUBTRACT: "-",
    MULTIPLY: "×",
    MULTIPLY_ALT: "*",
    IMPLIED_MULTIPLY: "imp-x",
    DIVIDE: "÷",
    DIVIDE_ALT: "/",
    EXPONENT: "^",
    // token types
    SYMBOL: "symbol",
    NUMBER: "number",
    ERROR: "error",
    // errors
    LEADING_ZERO: "leading zero",
    DIVIDE_ZERO: "divide by zero",
    UNKOWN: "unknown character",
    MISPLACED_DECIMAL: "misplaced decimal",
    NOT_NUMBER: "not a number",
    NO_PREFIX: "undefined prefix operation",
    NO_INFIX: "undefined infix operation",
    INCOMPLETE_EXPRESSION: "incomplete expression",
});
