import 'vue-router'
 
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean //依赖认证
    flatChildrenInMenu?: boolean //扁平化子菜单
    icon?: string  //图标
    hideInMenu?: boolean //不在菜单中显示
    locale?: string  //标题的国际化
    roles?: string[]  //可访问的角色
  }
}
