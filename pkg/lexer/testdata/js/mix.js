function normalize(str) { // @TEST_ANNOTATION fix bug in v8
  return str.split("\u0000");
}

/***/
