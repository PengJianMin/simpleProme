#/bin/bash
  
declare -a api
api[1]="/albums"
api[2]="/albums/delete"


declare -i number=$RANDOM
declare -i result=$number%2



# curl http://192.168.81.1:8080/albums     --include     --header "Content-Type: application/json"     --request "POST"     --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'

if [ $result -eq 0 ];then
        url="http://192.168.81.1:8080/albums"
        curl $url

        declare -i a=$number%6
        if [ $a -eq 0 ];then
                curl "$url/$RANDOM"
        fi
		
		declare -i b=$number%14
        if [ $b -eq 0 ];then
                   curl "$url" --include\
				--header "Content-Type: application/json"\
				--request "POST"\
				--data '{"id": "0","title": "The Modern Sound of Betty Carter(Verson:test_'"$RANDOM"')","artist": "test_'"$RANDOM-a"'","price": 49.99}'
        curl "$url" --include\
				--header "Content-Type: application/json"\
				--request "POST"\
				--data '{"id": "0","title": "The Modern Sound of Betty Carter(Verson:test_'"$RANDOM"')","artist": "test_'"$RANDOM-b"'","price": 49.99}'
        curl "$url" --include\
				--header "Content-Type: application/json"\
				--request "POST"\
				--data '{"id": "0","title": "The Modern Sound of Betty Carter(Verson:test_'"$RANDOM"')","artist": "test_'"$RANDOM-c"'","price": 49.99}'
        curl "$url" --include\
				--header "Content-Type: application/json"\
				--request "POST"\
				--data '{"id": "0","title": "The Modern Sound of Betty Carter(Verson:test_'"$RANDOM"')","artist": "test_'"$RANDOM-d"'","price": 49.99}'
        fi


else
        url="http://192.168.81.1:8080/albums/delete"
	curl "$url"
fi
