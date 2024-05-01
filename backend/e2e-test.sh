#!/usr/bin/env bash

set -euo pipefail
# if debug is needed, uncomment the line below
#set -x

# read optional host and port from command line args
host=${1:-localhost}
port=${2:-9000}
uri=http://$host:$port

echo "#############################################"
echo "Running e2e tests against $host:$port"
echo "#############################################"

echo "--- Test if server responds ---"
STATUSCODE=$(curl --silent --output /dev/stderr --write-out "%{http_code}" $uri)
if [ $STATUSCODE -ne 200 ]; then
  echo "Server is not running"
  exit 1
fi
echo ""

echo "--- User should not be authenticated ---"
STATUSCODE=$(curl --silent --output /dev/stderr --write-out "%{http_code}" $uri/authenticated)
if [ $STATUSCODE -ne  401 ]; then
	echo "Should have received 401 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- Creating a user ---"
tmp=$(mktemp)
STATUSCODE=$(curl -s -X POST --output $tmp --write-out "%{http_code}" $uri/create-user \
-d '{"name": "Alice", "email":"alice@gmail.com", "password":"123456789", "passwordAgain":"123456789"}') 
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
cat $tmp >&2
session=$(cat $tmp | jq -r '.session')
userid=$(cat $tmp | jq -r '.user.id')
echo ""

echo "--- User should be authenticated ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" -H "Cookie: session=$session" $uri/authenticated)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

tmp=$(mktemp)
echo "--- Quiz should be created ---"
STATUSCODE=$(curl -s -o $tmp --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	-d '{"name":"test-quiz", "description":"test-desc"}' \
	$uri/quizzes/create)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
cat $tmp >&2
quizid=$(cat $tmp | jq -r '.id')
echo ""

echo "--- Quizinfo should be returned ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" -H "Cookie: session=$session" $uri/quizzes/$quizid)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- Only quiz title should be modified ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	-X PATCH \
	-d '{"name":"test-quiz-modified"}' \
	$uri/quizzes/$quizid)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- Only quiz desc should be modified ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	-X PATCH \
	-d '{"description":"test-desc-modified"}' \
	$uri/quizzes/$quizid)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- Quiz desc and title should be modified ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	-X PATCH \
	-d '{"name":"test-name-mod2", "description":"test-desc-mod2"}' \
	$uri/quizzes/$quizid)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- All quizzes created by Alice should be returned ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	$uri/quizzes/user/$userid)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

tmp=$(mktemp)
echo "--- A multiple-choice question should be generated ---"
STATUSCODE=$(curl -s -o $tmp --write-out "%{http_code}" \
	-X POST \
	-H "Cookie: session=$session" \
	$uri/questions/multiple-choice \
	-d '{"quizId":"'$quizid'", "prompt":"Magyarország állam Közép-Európában, a Kárpát-medence közepén. 1989 óta parlamentáris köztársaság. Északról Szlovákia, északkeletről Ukrajna, keletről és délkeletről Románia, délről Szerbia, délnyugatról Horvátország és Szlovénia, nyugatról pedig Ausztria határolja."}')
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
multiple_q=$(cat $tmp | jq -r '.id')
cat $tmp >&2
echo ""

echo "--- A multiple-choice question should be returned ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	$uri/questions/multiple-choice/$multiple_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- A multiple-choice question should be modified ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-X PATCH \
	-H "Cookie: session=$session" \
	-d '{"quizId":"'$quizid'", "question":"Modified, whatever it was","answers":["CABD", "ABCD"], "CorrectAnswers":["A", "B"]}' \
	$uri/questions/multiple-choice/$multiple_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- A multiple-choice question should be deleted ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-X DELETE \
	-H "Cookie: session=$session" \
	$uri/questions/multiple-choice/$quizid/$multiple_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

tmp=$(mktemp)
echo "--- A single-choice question should be generated ---"
STATUSCODE=$(curl -s -o $tmp --write-out "%{http_code}" \
	-X POST \
	-H "Cookie: session=$session" \
	$uri/questions/single-choice \
	-d '{"quizId":"'$quizid'", "prompt":"Magyarország állam Közép-Európában, a Kárpát-medence közepén. 1989 óta parlamentáris köztársaság. Északról Szlovákia, északkeletről Ukrajna, keletről és délkeletről Románia, délről Szerbia, délnyugatról Horvátország és Szlovénia, nyugatról pedig Ausztria határolja."}')
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
single_q=$(cat $tmp | jq -r '.id')
cat $tmp >&2
echo ""

echo $single_q
echo "--- A single-choice question should be returned ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	$uri/questions/single-choice/$single_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- A single-choice question should be modified ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-X PATCH \
	-H "Cookie: session=$session" \
	-d '{"quizId":"'$quizid'", "question":"Modified, whatever it was","answers":["CABD", "ABCD"], "correctAnswer":"A"}' \
	$uri/questions/single-choice/$single_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- A single-choice question should be deleted ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-X DELETE \
	-H "Cookie: session=$session" \
	$uri/questions/single-choice/$quizid/$single_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

tmp=$(mktemp)
echo "--- A true-or-false question should be generated ---"
STATUSCODE=$(curl -s -o $tmp --write-out "%{http_code}" \
	-X POST \
	-H "Cookie: session=$session" \
	$uri/questions/true-or-false \
	-d '{"quizId":"'$quizid'", "prompt":"Magyarország állam Közép-Európában, a Kárpát-medence közepén. 1989 óta parlamentáris köztársaság. Északról Szlovákia, északkeletről Ukrajna, keletről és délkeletről Románia, délről Szerbia, délnyugatról Horvátország és Szlovénia, nyugatról pedig Ausztria határolja."}')
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
tf_q=$(cat $tmp | jq -r '.id')
cat $tmp >&2
echo ""

echo "--- A true-or-false question should be returned ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	$uri/questions/true-or-false/$tf_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- A true-or-false question should be modified ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-X PATCH \
	-H "Cookie: session=$session" \
	-d '{"QuizId":"'$quizid'", "Question":"Modified, whatever it was", "CorrectAnswer":true}' \
	$uri/questions/true-or-false/$tf_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- A true-or-false question should be deleted ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-X DELETE \
	-H "Cookie: session=$session" \
	$uri/questions/true-or-false/$quizid/$tf_q)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- Quiz should be deleted ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	-X DELETE \
	$uri/quizzes/$quizid)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""

echo "--- User should be deleted ---"
STATUSCODE=$(curl -s -o /dev/stderr --write-out "%{http_code}" \
	-H "Cookie: session=$session" \
	-X DELETE \
	$uri/delete-user/$userid)
if [ $STATUSCODE -ne  200 ]; then
	echo "Should have received 200 status code, got $STATUSCODE"
  exit 1
fi
echo ""
