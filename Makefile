
release:
	[[ $$( git rev-parse --abbrev-ref HEAD ) == "main" ]] # make sure we are on main
	git push origin main
	git tag $$( git tag | grep "^v" | sort --version-sort | tail -n 1 | awk -F. '{OFS="."; $$3 = $$3 + 1; print}' )
	git push --tags

goreleaser:
	goreleaser --snapshot --skip-publish --rm-dist
