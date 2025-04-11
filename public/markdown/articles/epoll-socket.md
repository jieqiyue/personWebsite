1. 可能有多个进程监听同一个socket，但是在三次握手之后，这些进程都会从io多路复用函数中返回。但是只有一个进程能够通过accept函数得到clientfd。其它的，如果是设置为阻塞的socket的话，其它没有从accept拿到clientfd的进程就会一直阻塞在那里。
2. 有些时候，虽然从select返回了，但是不一定能够通过read读取到数据。
3. 如果不将 clientfd 设置成非阻塞模式，那么一旦 epoll_wait 检测到读或者写事件返回后，接下来处理 clientfd 的读或写事件，如果对端因为 TCP 窗口太小，send 函数刚好不能将数据全部发送出去，将会造成阻塞，进而导致整个服务“卡住”。
4. 对于Redis来说，是将服务器端的listenfd设置为非阻塞。
在低版本的Redis中，仅仅是并没有将listenfd设置为非阻塞的。此时Redis服务端每次调用accept只会最多跟一个Redis客户端建立连接。可能会有性能问题。
而在高版本的Redis中，是有一个最大的连接数，也就是一次最多处理多少个客户端连接。也就是当io多路复用函数（epoll）返回报告listened上面有待处理的连接的时候，我们就会调用accept，这个时候来处理客户端连接。而具体要处理多少个，Redis中已经定义了最大的处理个数。

下面是Redis 6.0中服务器端在listenfd报告有链接到来的时候，怎样去调用accept函数来和客户端建立连接的函数。可以看到是有一个while循环来一次性拿多个链接的。
```c
// reids 6.0 版本
#define MAX_ACCEPTS_PER_CALL 1000

void acceptTcpHandler(aeEventLoop *el, int fd, void *privdata, int mask) {
    int cport, cfd, max = MAX_ACCEPTS_PER_CALL;
    char cip[NET_IP_STR_LEN];
    UNUSED(el);
    UNUSED(mask);
    UNUSED(privdata);

    while(max--) {
        cfd = anetTcpAccept(server.neterr, fd, cip, sizeof(cip), &cport);
        if (cfd == ANET_ERR) {
            if (errno != EWOULDBLOCK)
                serverLog(LL_WARNING,
                    "Accepting client connection: %s", server.neterr);
            return;
        }
        serverLog(LL_VERBOSE,"Accepted %s:%d", cip, cport);
        acceptCommonHandler(connCreateAcceptedSocket(cfd),0,cip);
    }
}

static int anetGenericAccept(char *err, int s, struct sockaddr *sa, socklen_t *len) {
    int fd;
    while(1) {
        fd = accept(s,sa,len);
        if (fd == -1) {
            if (errno == EINTR)  // 这个错误代表中断信号把accept打断了。此时可以continue进行下一次循环，来继续accept连接。
                continue;
            else {
                anetSetError(err, "accept: %s", strerror(errno));
                return ANET_ERR;
            }
        }
        break;   // 这里可以看出来，和这个是拿到clientfd立马返回的。
    }
    return fd;
}

int anetTcpAccept(char *err, int s, char *ip, size_t ip_len, int *port) {
    int fd;
    struct sockaddr_storage sa;
    socklen_t salen = sizeof(sa);
    if ((fd = anetGenericAccept(err,s,(struct sockaddr*)&sa,&salen)) == -1)
        return ANET_ERR;

    if (sa.ss_family == AF_INET) {
        struct sockaddr_in *s = (struct sockaddr_in *)&sa;
        if (ip) inet_ntop(AF_INET,(void*)&(s->sin_addr),ip,ip_len);
        if (port) *port = ntohs(s->sin_port);
    } else {
        struct sockaddr_in6 *s = (struct sockaddr_in6 *)&sa;
        if (ip) inet_ntop(AF_INET6,(void*)&(s->sin6_addr),ip,ip_len);
        if (port) *port = ntohs(s->sin6_port);
    }
    return fd;
}
```
在上面的代码中，有一个max变量，能够控制每次epoll报告listened有连接的时候，最多每次处理多少个客户端连接。

此时，在调用accept的时候，就不会阻塞了。我们可以根据错误码来判断当前listenfd是否真的没有客户端连接了。
## 问题

1. 那这里将listenfd设置为阻塞和非阻塞表现出来的行为有什么不同呢？
首先使用Redis1.3.6的源码进行调试，可以看到Redis1.3.6中在anetTcpServer中创建了服务端的listenfd之后，就直接去做事件循环了。并没有将listenfd设置为非阻塞。这个时候调用accept，将会一直阻塞在这里，直到有客户端连接过来。
![在这里插入图片描述](\images\art\8ed3bf74a0868c05f8bd97179e08de7f.png)然后我如果用一个客户端去连接的话，就会从accept中返回了。
![在这里插入图片描述](\images\art\f086e3807e32ae8c9d17868cde67e247.png)
而当我在anetTcpServer中创建了服务端的listenfd之后，手动添加一段代码，将这个listenfd设置为非阻塞呢？
![在这里插入图片描述](\images\art\07148e1a7a414de327f13ba8d2e195dc.png)
此时已经是非阻塞了。
此时再去执行accept，会直接返回。并且error number会设置为11.
![在这里插入图片描述](\images\art\c51a5a9d766af3624a2742eaad8c8450.png)
这样就能知道将listenfd设置为非阻塞的行为是什么了。


2. 这个和在一个时间循环中与多个客户端建立连接和每次时间循环只与一个客户端建立连接有什么区别呢？
下面通过实验，在Redis1.3.6中，如果同时有两个客户端与服务器端建立连接的时候，服务端是一个一个去建立连接的吗？
是的，由于Redis1.3.6在每一次时间循环中，只会accept一次，而一次accept最多只会返回一个可以连接的fd。



补充：


参考：
https://www.zhihu.com/question/37271342/answer/81550237

https://mp.weixin.qq.com/s?__biz=Mzk0MjUwNDE2OA==&mid=2247497596&idx=2&sn=5488887274663abd7b33d906fa6fcd9c&source=41#wechat_redirect