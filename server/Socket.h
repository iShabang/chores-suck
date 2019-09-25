#ifndef _SOCKET_H_
#define _SOCKET_H_

#include <string>

namespace fairmate { namespace server {

class Socket
{
public:
    Socket(std::string address);
    ~Socket();
    bool create(std::string address);
    bool close();
    int send(const std::string &buffer, int bytes);
    int recv(std::string &buffer, int bytes);

private:
    int m_sockFd;
    std::string m_address;

};

}// namespace server
}// namespace fairmate

#endif // _SOCKET_H_

