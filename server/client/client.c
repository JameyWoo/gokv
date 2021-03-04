/*************************************************************************
        > File Name: echo_client.c
        > Author: xjhznick
        > Mail: xjhznick@gmail.com
        > Created Time: 2015年03月17日 星期三 14时49分04秒
  > Description: Linux
 Socket网络编程--基于TCP的简单的回声Client端,向服务器端请求建立连接并接收服务器返回的回声字符串
 ************************************************************************/

#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

void error_handling(char* message);

#define BUF_SIZE 1024

int main(int argc, char* argv[]) {
    int server_sock;
    struct sockaddr_in server_addr;

    char amessage[BUF_SIZE];
    int str_len, recv_len, recv_cnt;

    if (3 != argc) {
        printf("Usage : %s <IP> <Port> \n", argv[0]);
        exit(1);
    }

    server_sock = socket(PF_INET, SOCK_STREAM, 0);
    if (-1 == server_sock) {
        error_handling("socket() error!");
        exit(1);
    }

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = inet_addr(argv[1]);
    server_addr.sin_port = htons(atoi(argv[2]));

    if (-1 == connect(server_sock, (struct sockaddr*)&server_addr,
                      sizeof(server_addr))) {
        error_handling("connect() error!");
    } else {
        puts("Connected......");
    }

    int t = 0;
    while (1) {
        char* command = (char*)malloc(128);
        int command_len = 0;
        if (t == 0) {
            printf("Input the db dir : ");
            fgets(command, BUF_SIZE, stdin);
            command_len = strlen(command) - 1;
            // ! 之所以要减去一个1是因为 输入的最后一个字符是 换行符\n
            printf("command_len: %d\n", command_len);
            t++;
        } else {
            printf("Input the Op : ");
            fgets(amessage, BUF_SIZE, stdin);
            command[0] = amessage[0] - '0';
            printf("Input the Key : ");
            fgets(amessage, BUF_SIZE, stdin);
            int key_len = strlen(amessage) - 1;
            // 要吐血了, 这里总是忘记 key_len 是会被除成0的!!!
            int key_len_2 = key_len;
            command_len = 5 + key_len_2;
            printf("key_len: %d\n", key_len);
            for (int i = 0; i < 4; i++) {
                command[i + 1] = key_len % 256;
                key_len /= 256;
            }
            for (int i = 0; i < key_len_2; i++) {
                command[5 + i] = amessage[i];
            }
            if (command[0] == 0) {  // PUT, 还要输入value
                printf("Input the Value : ");
                fgets(amessage, BUF_SIZE, stdin);
                int value_len = strlen(amessage) - 1;
                int value_len_2 = value_len;
                command_len += (4 + value_len);
                printf("value_len: %d\n", value_len);
                for (int i = 0; i < 4; i++) {
                    command[5 + key_len_2 + i] = value_len % 256;
                    value_len /= 256;
                }
                // fix bug: fuck啊, 又是简单的变量名搞错了!!! 搞成 key_len_2了,
                // 难怪!!!
                for (int i = 0; i < value_len_2; i++) {
                    command[9 + key_len_2 + i] = amessage[i];
                }
            }
        }

        int len = command_len;
        printf("len: %d\n", len);
        int xlen = len + 4;
        // 将 len 按小端序编码, 然后写入到一个新的 字符串中
        char* new_message = (char*)(malloc(xlen));
        for (int i = 0; i < 4; i++) {
            new_message[i] = len % 256;
            len /= 256;
        }
        for (int i = 4; i < xlen; i++) {
            new_message[i] = command[i - 4];
        }
        char* message = new_message;
        printf("command:\n");
        for (int i = 0; i < len; i++) {
            printf("%c", command[i]);
        }
        printf("\n");
        printf("message:\n");
        for (int i = 0; i < xlen; i++) {
            printf("%c", message[i]);
        }
        printf("\n");

        if (0 == strcmp(message, "q\n") || 0 == strcmp(message, "Q\n")) {
            break;
        }

        str_len = write(server_sock, message, xlen);
        printf("str_len: %d\n", str_len);

        recv_len = 0;
        // 这里需要先读取4字节, 然后得到长度, 之后再读取剩余的字节
        while (recv_len < 4) {
            recv_cnt =
                read(server_sock, &message[recv_len], BUF_SIZE - 1 - recv_len);
            if (-1 == recv_cnt) {
                error_handling("read() error!");
            }
            recv_len += recv_cnt;
        }
        int msg_len = 0;
        for (int i = 3; i >= 0; i--) {
            msg_len = msg_len * 256 + message[i];
        }
        while (recv_len < 4 + msg_len) {  //< 降低因異常情況陷入無限循環
            recv_cnt =
                read(server_sock, &message[recv_len], BUF_SIZE - 1 - recv_len);
            if (-1 == recv_cnt) {
                error_handling("read() error!");
            }
            recv_len += recv_cnt;
        }
        // 接收到全部字节之后
        message[recv_len] = 0;
        printf("Message from server : %s \n", &message[4]);
    }

    close(server_sock);

    return 0;
}

void error_handling(char* message) {
    fputs(message, stderr);
    fputc('\n', stderr);
    exit(1);
}
