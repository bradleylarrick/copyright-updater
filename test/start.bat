@echo off
setlocal enableextensions enabledelayedexpansion

set "IMAGE_NAME=%RAVENSTRIKE_DOCKER_IMAGE%"
if "%IMAGE_NAME%"=="" set "IMAGE_NAME=ravenstrike-python"

set "ROS_DISTRO_NAME=%ROS_DISTRO%"
if "%ROS_DISTRO_NAME%"=="" set "ROS_DISTRO_NAME=jazzy"

set "DASHBOARD_CMD=%RAVENSTRIKE_UI_CMD%"
if "%DASHBOARD_CMD%"=="" set "DASHBOARD_CMD=ros2 run ravenstrike_ui dashboard"
set "DOCKER_COLCON_ROOT=.docker_colcon"

pushd %~dp0 >nul
set "WORKSPACE_PATH=%cd%"

python --version >nul 2>&1
if errorlevel 1 (
    py -3 --version >nul 2>&1
    if errorlevel 1 (
        echo [error] Python 3 interpreter not found; install Python or add it to PATH.
        goto :error
    )
    set "PYTHON_EXE=py -3"
) else (
    set "PYTHON_EXE=python"
)

echo [preflight] validating declared DLL/library dependencies...
%PYTHON_EXE% scripts\preflight.py --warn-missing-env
if errorlevel 1 goto :error

docker info >nul 2>&1
if errorlevel 1 (
    echo [error] cannot connect to Docker; start Docker Desktop and retry.
    goto :error
)

echo [docker] building image %IMAGE_NAME% for ROS %ROS_DISTRO_NAME%...
docker build --build-arg ROS_DISTRO=%ROS_DISTRO_NAME% -t %IMAGE_NAME% .
if errorlevel 1 goto :error

set "RUN_CMD=mkdir -p %DOCKER_COLCON_ROOT% ^&^& source /opt/ros/%ROS_DISTRO_NAME%/setup.bash ^&^& colcon --log-base %DOCKER_COLCON_ROOT%/log build --symlink-install --build-base %DOCKER_COLCON_ROOT%/build --install-base %DOCKER_COLCON_ROOT%/install ^&^& source %DOCKER_COLCON_ROOT%/install/setup.bash ^&^& %DASHBOARD_CMD%"
set "DISPLAY_FLAG="
if defined DISPLAY set "DISPLAY_FLAG=-e DISPLAY=%DISPLAY%"

echo [docker] launching dashboard inside container...
docker run --rm -it -v "%WORKSPACE_PATH%:/workspace" -w /workspace -e ROS_DISTRO=%ROS_DISTRO_NAME% -e QT_QPA_PLATFORM=xcb -e QT_X11_NO_MITSHM=1 %DISPLAY_FLAG% %IMAGE_NAME% bash -lc "!RUN_CMD!"
set "EXIT_CODE=%ERRORLEVEL%"

popd >nul
exit /b %EXIT_CODE%

:error
popd >nul
exit /b 1
