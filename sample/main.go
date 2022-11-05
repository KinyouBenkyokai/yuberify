package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	issuerID   = "https://oxford.com/issuers/1" // This is a fake URL.
	issuerName = "The Marvelous University of Oxford"
)

func main() {
	// Part I: Create the issuer, the subject, and the verifier.
	issuer, err := CreateIssuer(issuerID, issuerName)
	if err != nil {
		panic(err)
	}

	subject, err := CreateSubject()
	if err != nil {
		panic(err)
	}

	verifier := CreateVerifier()

	// Part II: The Issuer issues credentials on the Subject.
	credentials, err := part2(issuer, subject)
	if err != nil {
		panic(err)
	}

	// Part III: The Verifier (any third party) can check the claim of the
	// Subject that it holds the credentials
	if err := part3(subject, verifier, credentials); err != nil {
		panic(err)
	}
}

func part2(issuer Issuer, subject Subject) (Credential, error) {
	// Step 1: Create a Subject and a claim to sign about this subject.
	// The claim is created jointly by the Subject and the Issuer. How they come
	// to agree on the claim to sign is out of scope here.
	claim := Claim{
		Age:            24,
		UniversityName: "Oxford",
		Degree:         "Bachelor of Science",
	}
	nicePrint(claim, "Claim")

	// Step 2: The Issuer signs the claim about this subject.
	id, err := subject.GetID()
	if err != nil {
		return Credential{}, err
	}
	credentials, err := issuer.SignCredential(claim, id)
	if err != nil {
		err = fmt.Errorf("Issuer couldn't sign credentials: %w", err)
		return credentials, err
	}

	nicePrint(credentials, "Credential")
	return credentials, err
}

func part3(subject Subject, verifier Verifier, credentials Credential) error {
	// Step 1: The verifier creates a challenge/nonce to be included in the
	// presentation which will be signed bby the subject.
	nonce, err := verifier.MakeNonce()
	if err != nil {
		return err
	}

	// Step 2: The subject creates the presentation and signs it.
	presentation, err := subject.SignPresentation(
		credentials,
		nonce,
	)
	if err != nil {
		return err
	}

	nicePrint(presentation, "Presentation")

	// Step 3: The verifier checks that the signature of the presentation is
	// correct.
	err = verifier.VerifiesPresentation(presentation)
	if err != nil {
		return fmt.Errorf("Verificiation failed: %w", err)
	}
	fmt.Println("\n!!! Verification succeeded !!!")
	return nil
}

func nicePrint(i interface{}, name string) {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")

	fmt.Printf("\n***** %s *****\n\n", name)
	e.Encode(i)
}
