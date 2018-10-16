build:
	docker build --no-cache -t dir-monitor:0.1 .
	docker tag dir-monitor:0.1 ai-image.jd.com/ceph/dir-monitor:0.1
	docker push ai-image.jd.com/ceph/dir-monitor:0.1
deploy-only-ht01:
	kubectl create ns dir-monitor || true
	kubectl delete -f deployment || true
	kubectl apply -f deployment/deployment-ht01.yaml
deploy-only-ht02:
	kubectl create ns dir-monitor || true
	kubectl delete -f deployment || true
	kubectl apply -f deployment/deployment-ht02.yaml
clean:
	kubectl delete ns dir-monitor || true
deploy-ht01: build deploy-only-ht01
deploy-ht02: build deploy-only-ht02
