#!/bin/bash

getOSName() {
    # shellcheck disable=SC2046
    os=$(trim $(cat /etc/os-release 2>/dev/null | grep ^ID= | awk -F= '{print $2}'))

    if [ "$os" = "" ]; then
        # shellcheck disable=SC2046
        os=$(trim $(lsb_release -i 2>/dev/null | awk -F: '{print $2}'))
    fi
    if [ ! "$os" = "" ]; then
        # shellcheck disable=SC2021
        os=$(echo "$os" | tr '[A-Z]' '[a-z]')
    fi

    echo "$os"
}

GetArchitecture() {
  _cputype="$(uname -m)"
  if [[ "$OSTYPE" == "linux-gnu" ]]; then
          set -e
          if [[ $(whoami) == "root" ]]; then
                  MAKE_ME_ROOT=
          else
                  MAKE_ME_ROOT=sudo
          fi
          echo "echo Arch Linux detected."
          os="linux"
          deptool="ldd"
          osname=getOSName
          goAchive="go1.17.linux-amd64.tar.gz"
          libArray=("/usr/local/lib/librocksdb.6.27.dylib" "/usr/local/lib/librocksdb.6.dylib" "/usr/local/lib/librocksdb.dylib")
          rocksdbdep=("/usr/local/lib/pkconfig/rocksdb.pc" "/usr/local/include/rocksdb/" "/usr/local/lib/librocksdb.so*")
          case ${osname} in
              ubuntu)
                  echo "detected ubuntu ..."
                  dynamicLink="TRUE"
                  ;;
              centos)
                  echo "detected centos ..."
                  dynamicLink="TRUE"
                  ;;
              alpine)
                  echo "detected alpine ..."
                  dynamicLink="FALSE"
                  ;;
              *)
                  echo unknow os $OS, exit!
                  return
                  ;;
          esac
  elif [[ "$OSTYPE" == "darwin"* ]]; then
          set -e
          echo "Mac OS (Darwin) detected."
          os="darwin"
          deptool="otool -L"
          libArray=("/usr/local/lib/librocksdb.6.27.dylib" "/usr/local/lib/librocksdb.6.dylib" "/usr/local/lib/librocksdb.dylib")
          rocksdbdep=("" "" "")
          dynamicLink="TRUE"
          if [ "$_cputype" == "arm64" ]
          then
            goAchive="go1.17.darwin-arm64.tar.gz"
          else
            goAchive="go1.17.darwin-amd64.tar.gz"
          fi

  else
          echo "Unknown operating system."
          echo "This OS is not supported with this script at present. Sorry."
          exit 1
  fi
}

download() {
  rm -rf "$HOME"/.exchain/src
  mkdir -p "$HOME"/.exchain/src
  tag=`wget -qO- -t1 -T2 "https://api.github.com/repos/okex/exchain/releases/latest" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
  wget "https://github.com/okex/exchain/archive/refs/tags/${tag}.tar.gz" -O "$HOME"/.exchain/src/exchain.tar.gz
  ver=$(echo $tag| sed 's/v//g')
  cd "$HOME"/.exchain/src && tar zxvf exchain.tar.gz &&  cd exchain-"$ver"
}

function checkgoversion { echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4); }'; }

installRocksdb() {
  echo "install rocksdb...."
  make rocksdb
  echo "rocksdb install completed"
}

uninstallRocksdb() {
   # shellcheck disable=SC2068
   for lib in ${rocksdbdep[@]}
    do
      echo "rm lib ${lib}"
      rm -rf $lib
    done
  end
  echo "uninstallRocksdb ..."
}

installgo() {
  echo "install go ..."
  wget "https://golang.google.cn/dl/${goAchive}"
  tar -zxvf ${goAchive} -C /usr/local/
  rm go1.17.linux-amd64.tar.gz
  cd ~
  if [[ -f ".bashrc" ]]; then
      echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
      source ~/.bashrc
  fi
  if [[ -f ".zshrc" ]]; then
      echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.zshrc
      source ~/.zshrc
  fi

  /usr/local/go/bin/go env -w GOPROXY=https://goproxy.cn,direct
  /usr/local/go/bin/go env -w GO111MODULE="auto"
  echo "install go completed ..."
}

installWasmLib() {
  wget "https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0/libwasmvm_muslc.x86_64.a" -O /lib/libwasmvm_muslc.x86_64.a
  cp /lib/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.a
}

Prepare() {
  echo "check go version"
  v=$(go version | { read _ _ v _; echo ${v#go}; })
  # shellcheck disable=SC2046
  if [ $(checkgoversion "$v") -ge $(checkgoversion "1.17") ]
  then
    echo "$v"
    echo "should not install go"
  else
    echo "should install go version above 1.17"
    installgo
  fi

  echo "check jcmalloc ..."
  # shellcheck disable=SC2068
  for lib in ${libArray[@]}
  do
  echo "check lib ${lib}"
  if [ -f "${lib}" ]; then
    has=$($deptool -L "${lib}" |grep jemalloc |wc -l)
    if [ "${has}" -gt 0 ]; then
          uninstallRocksdb
          installRocksdb
          return
    fi
  fi
  done

  echo "check rocksdb lib version"
  # shellcheck disable=SC2068
  for lib in ${libArray[@]}
  do
    if [ ! -f "${lib}" ]; then
      uninstallRocksdb
      installRocksdb
      return
    fi
  done
  echo "Prepare completed ...."
}

InstallExchain() {
  echo "InstallExchain...."
  download
  #if alpine add LINK_STATICALLY=true
  echo "compile exchain...."
  if [ "$dynamicLink" == "TRUE" ]; then
    make mainnet WITH_ROCKSDB=true
  else
    installWasmLib
    make mainnet WITH_ROCKSDB=true LINK_STATICALLY=true
  fi

  echo "InstallExchain completed"
}

GetArchitecture
Prepare
InstallExchain