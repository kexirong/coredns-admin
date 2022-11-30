import { Space, Form, Button, Input, Card, Layout, Menu, Dropdown, Tree, Table, Scrollbar, InputNumber, Select, Popconfirm, Alert } from '@arco-design/web-vue'
import { FormItem, CardMeta, LayoutContent, LayoutHeader, LayoutSider, Row, Col, MenuItem, SubMenu } from '@arco-design/web-vue/es'
import Doption from '@arco-design/web-vue/es/dropdown/dropdown-option'

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
    app.component('ADoption', Doption)
    app.component('ARow', Row)
    app.component('ACol', Col)
    app.component('ATree', Tree)
    app.component('ATable', Table)
    app.component('AScrollbar', Scrollbar)
    app.component('AInputNumber', InputNumber)
    app.component('ASelect', Select)
    app.component('APopconfirm', Popconfirm)
    app.component('AAlert', Alert)
}
