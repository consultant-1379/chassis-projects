set(GTEST_SHUFFLE 1)

add_executable(
    ut
    main.cpp
    test.cpp
)

set_target_properties(
    ut
    PROPERTIES
    OUTPUT_NAME ut
)

set_target_properties(
    ut
    PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY ${PROJECT_SOURCE_DIR}/ut/
)

target_compile_options( 
    ut
    PRIVATE ${GCC_COVERAGE_COMPILE_FLAGS}
)

target_include_directories(
    ut 
    PUBLIC 
    ../vendor/googletest/googletest/include 
    ${CMAKE_SOURCE_DIR}/include
)

target_link_libraries(ut PUBLIC gtest)
