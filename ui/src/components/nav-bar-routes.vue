<script>
// Allows stubbing BaseLink in unit tests
const BaseLink = 'BaseLink'

export default {
  // Functional components are stateless, meaning they can't
  // have data, computed properties, etc and they have no
  // `this` context.
  functional: true,
  props: {
    routes: {
      type: Array,
      required: true,
    },
  },
  // Render functions are an alternative to templates
  render(h, { props, $style = {} }) {
    function getRouteTitle(route) {
      return typeof route.title === 'function' ? route.title() : route.title
    }

    function getRouteBadge(route) {
      if (!route.badge) {
        return false
      }
      return typeof route.badge === 'function' ? route.badge() : route.badge
    }

    // Functional components are the only components allowed
    // to return an array of children, rather than a single
    // root node.
    return props.routes.map((route) => {
      if (getRouteBadge(route) > 0) {
        return (
          <BaseLink
            tag='b-navbar-item'
            key={route.name}
            to={route}
            exact-active-class={$style.active}
          >
            <a>{getRouteTitle(route)}</a>
            <b-tag rounded type='is-danger is-light'>
              {getRouteBadge(route)}
            </b-tag>
          </BaseLink>
        )
      } else {
        return (
          <BaseLink
            tag='b-navbar-item'
            key={route.name}
            to={route}
            exact-active-class={$style.active}
          >
            <a>{getRouteTitle(route)}</a>
          </BaseLink>
        )
      }
    })
  },
}
</script>
