#ifndef _PASSIVE_SOCKET_H_
#define _PASSIVE_SOCKET_H_

#include <string>
#include "utils/Thread.h"

namespace fairmate { namespace server {

class PassiveSocket {
public:
    PassiveSocket(std::string address, const long long &queueSize);
    ~PassiveSocket();
    bool create(std::string address);
    bool close();
    int listen();
private:
    int m_sockFd;
    std::string m_address;

private:
    class SocketThread : public utils::Thread {
    public:
        SocketThread(long long queueSize);
        void operator()() override;
        void shutdown() override;
        void enable();
        void setSocket(int fd);
    private:
        bool m_shutdown;
        long long m_queueSize;
        int m_sockFd;
    };

private:
    SocketThread m_thread;

};

}// namespace server
}// namespace fairmate

#endif // _PASSIVE_SOCKET_H_

