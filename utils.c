#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *copy_string(const char *s) {
  int len = 0;
  while (s[len] != '\0') {
    len++;
  }
  char *new_str = (char *)malloc((len + 1) * sizeof(char));
  for (int i = 0; i <= len; i++) {
    new_str[i] = s[i];
  }
  return new_str;
}

int removeChar(char *str, char charToRemove) {
  int i, j;
  int len = strlen(str);
  int removed = false;

  // Iterate through the string
  for (i = 0, j = 0; i < len; i++) {
    // Copy only if current char is not the one to remove
    if (str[i] == charToRemove) {
      removed = true;
    } else {
      str[j] = str[i];
      j++;
    }
  }
  // Add null terminator at the end
  str[j] = '\0';

  return removed;
}
