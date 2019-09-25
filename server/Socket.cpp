#include <sys/socket.h>

#include "Socket.h"

namespace fairmate { namespace server {

Socket::Socket(std::string address):
    m_sockFd(-1),
    m_address(address)
{}

bool Socket::create(std::string address)
{
    bool ret = false;
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

    return ret;
}

} // namespace server
} // namespace fairmate
