import { Space, Form, Button, Input, Card, Layout, Menu, Dropdown } from '@arco-design/web-vue'
import { FormItem, CardMeta, LayoutContent, LayoutHeader, LayoutSider, MenuItem, SubMenu, } from '@arco-design/web-vue/es'
import DOption from '@arco-design/web-vue/es/dropdown/dropdown-option'

import { App } from 'vue'
export default (app: App<Element>) => {
    app.component('AForm', Form)
    app.component('ASpace', Space)
    app.component('AButton', Button)
    app.component('AInput', Input)
    app.component('AFormItem', FormItem)
    app.component('ACard', Card)
    app.component('ACardMeta', CardMeta)
    app.component('ALayout', Layout)
    app.component('ALayoutContent', LayoutContent)
    app.component('ALayoutHeader', LayoutHeader)
    app.component('ALayoutSider', LayoutSider)
    app.component('AMenu', Menu)
    app.component('AMenuItem', MenuItem)
    app.component('ASubMenu', SubMenu)
    app.component('ADropdown', Dropdown)
    app.component('ADOption', DOption)
}
