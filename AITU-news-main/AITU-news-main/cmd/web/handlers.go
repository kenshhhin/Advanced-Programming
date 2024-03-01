package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"xximsz.net/snippetbox/pkg/forms"
	"xximsz.net/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form data.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		// If the form data is not valid, render the create snippet form again
		// with validation errors.
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// Insert the snippet into storage.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		// If there's an error during insertion, return a server error.
		app.serverError(w, err)
		return
	}

	// Set a flash message indicating success.
	app.session.Put(r, "flash", "Snippet successfully created!")

	// Redirect the user to the newly created snippet's page.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil)})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 8)

	// If the form is not valid, redisplay the signup form with validation errors.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	// Try to create a new user record in the database.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		// If the email already exists, add an error message to the form and re-display it.
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	// Add a confirmation flash message to the session confirming successful signup.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	// Redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil)})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		// If the credentials are invalid, add a generic error message to the form and re-display the login page.
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	// Add the ID of the current user to the session to indicate they are logged in.
	app.session.Put(r, "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is 'logged out'.
	app.session.Remove(r, "authenticatedUserID")

	// Add a flash message to the session to confirm to the user that they've been logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")

	// Redirect the user to the home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type Comment struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	NewsID int    `json:"news_id"`
	Text   string `json:"text"`
}
