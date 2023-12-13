// filename: printenv.c
#include <stdio.h>
#include <stdlib.h>

int main() {
    const char* val = getenv("TEST_VAR");
    if (val) {
        printf("%s\n", val);
    } else {
        printf("TEST_VAR not set\n");
    }
    return 0;
}
