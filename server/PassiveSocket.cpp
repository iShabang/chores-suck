#include <sys/socket.h>

#include "PassiveSocket.h"

namespace fairmate { namespace server {

PassiveSocket::PassiveSocket(std::string address, const long long &queueSize):
    m_sockFd(-1),
    m_address(address),
    m_thread(queueSize)
{}

bool PassiveSocket::create(std::string address)
{
    m_sockFd = socket(AF_INET, SOCK_STREAM, 0);
    if (!m_sockFd)
    {
        return false;
    }
    struct sockaddr adr; 
    adr.sa_family = AF_INET;
    const char *buff = address.c_str();
    std::copy(buff, buff+address.size()-1, adr.sa_data);

    if (!bind(m_sockFd, &adr, sizeof(adr)))
    {
        return false;
    }

    return true;
}

PassiveSocket::SocketThread::SocketThread(long long queueSize):
    m_shutdown(false),
    m_queueSize(queueSize)
{
}

void PassiveSocket::SocketThread::setSocket(int fd)
{
    m_sockFd = fd;
}

void PassiveSocket::SocketThread::shutdown()
{
    m_shutdown = true;
}

void PassiveSocket::SocketThread::enable()
{
    m_shutdown = false;
}

void PassiveSocket::SocketThread::operator()()
{
    while(!m_shutdown){
        if (!::listen(m_sockFd,m_queueSize))
        {
            break;
        }
        struct sockaddr adr; 
        socklen_t size;
        int newFd = ::accept(m_sockFd, &adr, &size);
    }
}

} // namespace server
} // namespace fairmate
