HEADERS += headers/ItemInterface.h \
           headers/ItemHandler.h \
           server/Socket.h \ 
           server/PassiveSocket.h \ 
           utils/Thread.h \ 

SOURCES += main.cpp \
           sources/ItemHandler.cpp \
           server/PassiveSocket.cpp \
           utils/Thread.cpp

TARGET = fairmate
