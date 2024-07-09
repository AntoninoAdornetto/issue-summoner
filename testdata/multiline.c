#include <stdio.h>
#include <stdlib.h>

struct Person {
  int /* @TEST_TODO inline comment #1 */ age;
  char *name /* @TEST_TODO inline comment #2 */;
};

int main(int argc, char *argv[]) {
  char *no_comment_token =
      "This should not end up in the comment token / /* */ @TEST_TODO";
  printf("%s\n", no_comment_token);
  return 0;
}

/*
 * @TEST_TODO drop a star if you know about this code wars challenge
 * Digital Cypher assigns to each letter of the alphabet unique number.
 * Instead of letters in encrypted word we write the corresponding number
 * Then we add to each obtained digit consecutive digits from the key
 * */
char *decode(const unsigned char *code, size_t n, unsigned key) {
  char *msg = calloc(n + 1, 1);
  char buf[64];
  int key_len = sprintf(buf, "%d", key);

  for (size_t i = 0; i < n; i++) {
    msg[i] = code[i] - buf[i % key_len] + '0' + 'a' - 1;
  }

  return msg;
}

// This comment should not be parsed since it does not have an annotation
