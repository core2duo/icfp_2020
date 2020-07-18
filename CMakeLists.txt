cmake_minimum_required(VERSION 3.6)
project(icpf_2020)

SET(CMAKE_C_COMPILER /usr/local/bin/gcc-9)
SET(CMAKE_CXX_COMPILER /usr/local/bin/g++-9)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fopenmp")

#set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -DDEBUG -O0 -fsanitize=address")
set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -DDEBUG -O0")
set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} -O3")

FILE(GLOB SOURCE_FILES app/main.cpp app/*.h app/appliers/*.h)



add_executable(icpf_2020 ${SOURCE_FILES})