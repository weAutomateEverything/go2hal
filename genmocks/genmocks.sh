#!/bin/bash
for LQN in $(ls  ../*/*mock.go );do
  DIR=$(echo ${LQN}| awk -F/ '{print $2}')
  SRC=$(echo ${LQN}| sed  's/\_mock//' )
  mockgen -source ${SRC} -package=${DIR} -destination=${LQN}

done
