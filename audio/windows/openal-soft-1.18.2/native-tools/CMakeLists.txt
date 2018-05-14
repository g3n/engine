cmake_minimum_required(VERSION 3.0.2)
project(native-tools)
add_executable(bin2h bin2h.c)
# Enforce no dressing for executable names, so the main script can find it
set_target_properties(bin2h PROPERTIES OUTPUT_NAME bin2h)
# Avoid configuration-dependent subdirectories while building with Visual Studio
set_target_properties(bin2h PROPERTIES RUNTIME_OUTPUT_DIRECTORY_DEBUG "${CMAKE_BINARY_DIR}")
set_target_properties(bin2h PROPERTIES RUNTIME_OUTPUT_DIRECTORY_RELEASE "${CMAKE_BINARY_DIR}")
