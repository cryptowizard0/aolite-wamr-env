#include <stdio.h>

__attribute__((export_name("print_string")))
void print_string(const char* str) {
    printf("Received string from Go: %s\n", str);
}

int main() {
    return 0;
}