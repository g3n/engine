project(alsoft-config)

option(ALSOFT_NO_QT5 "Use Qt4 instead of Qt5 for alsoft-config" FALSE)

include_directories("${alsoft-config_BINARY_DIR}")

set(alsoft-config_SRCS  main.cpp
                        mainwindow.cpp
)
set(alsoft-config_UIS  mainwindow.ui)
set(alsoft-config_MOCS  mainwindow.h)

find_package(Qt5Widgets)
if(Qt5Widgets_FOUND AND NOT ALSOFT_NO_QT5)
    qt5_wrap_ui(UIS ${alsoft-config_UIS})

    qt5_wrap_cpp(MOCS ${alsoft-config_MOCS})

    add_executable(alsoft-config ${alsoft-config_SRCS} ${UIS} ${RSCS} ${TRS} ${MOCS})
    target_link_libraries(alsoft-config Qt5::Widgets)
    set_property(TARGET alsoft-config APPEND PROPERTY COMPILE_FLAGS ${EXTRA_CFLAGS})
    set_target_properties(alsoft-config PROPERTIES RUNTIME_OUTPUT_DIRECTORY ${OpenAL_BINARY_DIR})
    if(TARGET build_version)
        add_dependencies(alsoft-config build_version)
    endif()

    install(TARGETS alsoft-config
            RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
            LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
            ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
    )
else()
    # Need Qt 4.8.0 or newer for the iconset theme attribute to work
    find_package(Qt4 4.8.0 COMPONENTS QtCore QtGui)
    if(QT4_FOUND)
        include(${QT_USE_FILE})

        qt4_wrap_ui(UIS ${alsoft-config_UIS})

        qt4_wrap_cpp(MOCS ${alsoft-config_MOCS})

        add_executable(alsoft-config ${alsoft-config_SRCS} ${UIS} ${RSCS} ${TRS} ${MOCS})
        target_link_libraries(alsoft-config ${QT_LIBRARIES})
        set_property(TARGET alsoft-config APPEND PROPERTY COMPILE_FLAGS ${EXTRA_CFLAGS})
        set_target_properties(alsoft-config PROPERTIES RUNTIME_OUTPUT_DIRECTORY ${OpenAL_BINARY_DIR})
        if(TARGET build_version)
            add_dependencies(alsoft-config build_version)
        endif()

        install(TARGETS alsoft-config
                RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
                LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
                ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
        )
    endif()
endif()
