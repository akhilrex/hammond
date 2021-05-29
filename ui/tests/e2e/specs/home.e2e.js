describe('Home Page', () => {
  it('has the correct title and heading', () => {
    cy.visit('/')
    cy.title().should('equal', 'Home | Hammond')
    cy.contains('h1', 'Home Page')
  })
})
