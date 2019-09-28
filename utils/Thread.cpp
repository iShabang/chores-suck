#include "Thread.h"

namespace fairmate { namespace utils {

Thread::Thread()
{
}

Thread::~Thread()
{
}

void *Thread::threadFcn(void *t)
{
    Thread *thread = static_cast<Thread *>(t);
    (*thread)();
    thread->join();
    return t;
}

void Thread::join()
{
    void *status;
    pthread_join(m_thread,&status);
}

void Thread::start()
{
    pthread_attr_t attr;
    pthread_attr_init(&attr);
    pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_JOINABLE);
    m_thread = pthread_create(&m_thread, &attr, threadFcn, static_cast<void*>(this));
}

} // namespace utils
} // namespace fairmate
