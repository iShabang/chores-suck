#ifndef _SOCKET_H_
#define _SOCKET_H_

#include <string>

#include "Socket.h"

namespace fairmate { namespace server {

class PassiveSocket : public Socket
{
public:
    PassiveSocket(std::string address);
    ~PassiveSocket();
    bool create(std::string address);
    bool close();
    int listen();
private:
    int m_sockFd;
    std::string m_address;
};

}// namespace server
}// namespace fairmate

#endif // _SOCKET_H_

