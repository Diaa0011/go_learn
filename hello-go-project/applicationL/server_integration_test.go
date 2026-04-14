package main

// func TestRecordingWinsAndRetrievingThem(t *testing.T) {
// 	store := NewInMemoryPlayerStore()
// 	server := PlayerServer{store}
// 	player := "Pepper"

// 	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
// 	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
// 	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

// 	response := httptest.NewRecorder()

// 	server.ServeHTTP(response, newGetScoreRequest(player))
// 	assertStatus(t, response.Code, http.StatusOK)

// 	assertResponseBody(t, response.Body.String(), "3")

// }
// func TestRecordingWinsAndRetrievingThem(t *testing.T) {
// 	store := NewInMemoryPlayerStore()
// 	server := PlayerServer{store}
// 	player := "Pepper"

// 	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
// 	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
// 	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

// 	response := httptest.NewRecorder()

// 	server.ServeHTTP(response, newGetScoreRequest(player))
// 	assertStatus(t, response.Code, http.StatusOK)

// 	assertResponseBody(t, response.Body.String(), "3")

// }
// }
