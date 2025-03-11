#include <stdio.h>
#include <stdlib.h>
#include <unistd.h> // 用于 write()，测试标准输出
#include <string.h>  // 添加 strlen 声明

int main() {
    char *buf;

    printf("Hello world!\n");

    buf = malloc(1024);
    if (!buf) {
        printf("malloc buf failed\n");
        return -1;
    }

    printf("buf ptr: %p\n", buf);

    sprintf(buf, "%s", "1234\n");
    printf("buf: %s", buf);

    free(buf);
    return 0;
}


