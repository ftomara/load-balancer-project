algos=("RoundRobin" "WRoundRobin" "Random" "LeastConnections" "Hash")
for id in 1 2 3 4 5; do
    echo "Running algorithm $id: ${algos[$id]}"
    ab -n 10000 -c 200 \
       -e "results/algo${id}_percentiles.csv" \
       "http://localhost:8080/calc?n=5&id=${id}" \
       | tee "results/algo${id}_summary.txt"
    echo "Done with algorithm $id"
    sleep 5 
done
echo "All benchmarks complete!"