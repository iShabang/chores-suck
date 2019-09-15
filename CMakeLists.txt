cmake_minimum_required(VERSION 3.1.0)

project(FairMate)

set(CMAKE_INCLUDE_CURRENT_DIR ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTOUIC ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

#include_directories(spent-core/Headers)

set(core_SRCS
    #spent-core/Sources/DatabaseManager.cpp
    #spent-core/Sources/Transaction.cpp
    #spent-core/Sources/TransactionDao.cpp
    #spent-core/Sources/TransactionModel.cpp
    main.cpp
)

add_executable(FairMate ${core_SRCS})

#target_link_libraries(Fairmate Qt5::Core Qt5::Sql)
