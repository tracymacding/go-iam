RUNBIN="${BASH_SOURCE-$0}"
RUNBIN="$(dirname "${RUNBIN}")"
BINDIR="$(cd "${RUNBIN}"; pwd)"

LOGDIR="logs"

while getopts "l:m" opt; do
    case $opt in
      l)
        LOGDIR=$OPTARG
        ;;
      m)
        MONGOADDRESS=$OPTARG
        ;;
    esac
done

mkdir -p $LOGDIR

cd $BINDIR && nohup ./galaxy-s3-gateway -log_dir=$LOGDIR -logtostderr=false -mongodb_addr=$MONGOADDRESS > nohup.out 2>&1 &

PID=$!
sleep 1
echo "start galaxy-s3-gateway finished, PID=$PID"
echo "checking if $PID is running..."
sleep 2
kill -0 $PID > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "$PID is running, start galaxy-s3-gateway success."
	exit 0
else
	echo "start galaxy-s3-gateway failed."
	exit 1
fi
