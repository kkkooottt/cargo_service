create_postgres:
	docker run --name testpg -p 5432:5432 -e POSTGRES_USER=kkkooottt -e POSTGRES_PASSWORD=secretpassword -d postgres

start_postgres:
	docker start testpg
stop_postgres:
	docker stop testpg