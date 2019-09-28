#ifndef _FAIRMATE_THREAD_H_
#define _FAIRMATE_THREAD_H_

#include <pthread.h>

namespace fairmate { namespace utils {
class Thread
{
public:
    Thread();
    ~Thread();
    void start();
    static void *threadFcn(void *t);
    void join();
    virtual void operator()() = 0;
    virtual void shutdown() = 0;

private:
    pthread_t m_thread;
    int m_id;

};
} // namespace utils
} // namespace fairmate

#endif // _FAIRMATE_THREAD_H_
