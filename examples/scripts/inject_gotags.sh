#!/bin/bash

script_dir=$(
    cd $(dirname $0)
    pwd
)                                  # 脚本路径
project_dir=$(dirname $script_dir) # 项目路径
out_dir=${project_dir}/gen         # 生成代码路径

files=$(find ${out_dir} -type f -name '*.pb.go')
for file in $files; do
    protoc-go-inject-tag -input=$file
done
